package kibana

import (
	"github.com/parnurzeal/gorequest"
)

type HttpAgent struct {
	client      *gorequest.SuperAgent
	authHandler AuthenticationHandler
	config      *Config
}

type AuthenticationHandler interface {
	Initialize(agent *gorequest.SuperAgent) error
	ChangeAccount(accountId string, agent *HttpAgent) error
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
		authHandler: authHandler,
		config:      config,
	}
}

func (authClient *HttpAgent) Auth(handler AuthenticationHandler) *HttpAgent {
	authClient.authHandler = handler
	return authClient
}

func (authClient *HttpAgent) Get(targetUrl string) *HttpAgent {
	agent := authClient.clone()
	agent.client.Get(targetUrl)
	return agent
}

func (authClient *HttpAgent) Delete(targetUrl string) *HttpAgent {
	agent := authClient.clone()
	agent.client.Delete(targetUrl)
	return agent
}

func (authClient *HttpAgent) Put(targetUrl string) *HttpAgent {
	agent := authClient.clone()
	agent.client.Put(targetUrl)
	return agent
}

func (authClient *HttpAgent) Post(targetUrl string) *HttpAgent {
	agent := authClient.clone()
	agent.client.Post(targetUrl)
	return agent
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

func (auth *BasicAuthenticationHandler) ChangeAccount(accountId string, agent *HttpAgent) error {
	return nil
}

func (auth *NoAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	return nil
}

func (auth *NoAuthenticationHandler) ChangeAccount(accountId string, agent *HttpAgent) error {
	return nil
}

func (authClient *HttpAgent) clone() *HttpAgent {
	return &HttpAgent{authHandler: authClient.authHandler, client: authClient.createSuperAgent(), config: authClient.config}
}

func (authClient *HttpAgent) createSuperAgent() *gorequest.SuperAgent {
	superAgent := gorequest.New()
	superAgent.Debug = authClient.config.Debug

	return superAgent
}
