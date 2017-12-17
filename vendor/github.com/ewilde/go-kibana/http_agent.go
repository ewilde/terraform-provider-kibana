package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type HttpAgent struct {
	client      *gorequest.SuperAgent
	authHandler AuthenticationHandler
}

type AuthenticationHandler interface {
	Initialize(agent *gorequest.SuperAgent) error
}

type NoAuthenticationHandler struct {
}

type BasicAuthenticationHandler struct {
	userName string
	password string
}

type LogzAuthenticationHandler struct {
	Auth0Uri     string
	LogzUri      string
	UserName     string
	Password     string
	ClientId     string
	sessionToken string
}

type Auth0Response struct {
	IdTokens    string `json:"id_token"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewHttpAgent(config *Config, authHandler AuthenticationHandler) *HttpAgent {
	return &HttpAgent{
		client:      gorequest.New(),
		authHandler: authHandler,
	}
}

func (authClient *HttpAgent) Auth(handler AuthenticationHandler) *HttpAgent {
	authClient.authHandler = handler
	return authClient
}

func (authClient *HttpAgent) Get(targetUrl string) *HttpAgent {
	authClient.client.Get(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Delete(targetUrl string) *HttpAgent {
	authClient.client.Delete(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Put(targetUrl string) *HttpAgent {
	authClient.client.Put(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Post(targetUrl string) *HttpAgent {
	authClient.client.Post(targetUrl)
	return authClient
}

func (authClient *HttpAgent) Query(content interface{}) *HttpAgent {
	authClient.client.Query(content)
	return authClient
}

func (authClient *HttpAgent) Set(param string, value string) *HttpAgent {
	authClient.client.Set(param, value)
	return authClient
}

func (authClient *HttpAgent) Send(content interface{}) *HttpAgent {
	authClient.client.Send(content)
	return authClient
}

func (authClient *HttpAgent) End(callback ...func(response gorequest.Response, body string, errs []error)) (gorequest.Response, string, []error) {
	if err := authClient.authHandler.Initialize(authClient.client); err != nil {
		return nil, "", []error{err}
	}

	return authClient.client.End(callback...)
}

func NewBasicAuthentication(userName string, password string) *BasicAuthenticationHandler {
	return &BasicAuthenticationHandler{userName: userName, password: password}
}

func (auth *BasicAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	agent.SetBasicAuth(auth.userName, auth.password)
	return nil
}

func (auth *NoAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	return nil
}

func (auth *LogzAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	if auth.sessionToken != "" {
		auth.setLogzHeaders(agent)
		return nil
	}

	request := gorequest.New()
	response, body, errs := request.Post(fmt.Sprintf("%s/oauth/ro", auth.Auth0Uri)).
		Set("kbn-version", DefaultKibanaVersion553).
		Set("Content-Type", "application/x-www-form-urlencoded").
		Type("form").
		Send(fmt.Sprintf(`{
  "scope": "openid email connection",
  "response_type": "code",
  "connection": "Username-Password-Authentication",
  "username": "%s",
  "password": "%s",
  "grant_type": "password",
  "client_id": "%s"
}`, auth.UserName, auth.Password, auth.ClientId)).
		End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	authResponse := &Auth0Response{}
	if err := json.Unmarshal([]byte(body), authResponse); err != nil {
		return err
	}

	response, body, errs = request.Post(fmt.Sprintf("%s/login/jwt", auth.LogzUri)).
		Send(fmt.Sprintf(`{
  "jwt": "%s"
}`, authResponse.IdTokens)).
		End()

	jwtResponse := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &jwtResponse); err != nil {
		return err
	}

	auth.sessionToken = jwtResponse["sessionToken"].(string)
	auth.setLogzHeaders(agent)
	return nil
}

func (auth *LogzAuthenticationHandler) setLogzHeaders(agent *gorequest.SuperAgent) *gorequest.SuperAgent {
	return agent.
		Set("Content-Type", "application/json").
		Set("x-auth-token", auth.sessionToken)

}
