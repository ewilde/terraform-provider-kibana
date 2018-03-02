package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

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

func (auth *LogzAuthenticationHandler) ChangeAccount(accountId string, agent *HttpAgent) error {
	response, body, errs := agent.Get(fmt.Sprintf("%s/user/session/replace/%s", auth.LogzUri, accountId)).End()
	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	responseMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &responseMap); err != nil {
		return err
	}

	auth.sessionToken = responseMap["sessionToken"].(string)
	return nil
}

func (auth *LogzAuthenticationHandler) setLogzHeaders(agent *gorequest.SuperAgent) *gorequest.SuperAgent {
	return agent.
		Set("Content-Type", "application/json").
		Set("x-auth-token", auth.sessionToken)

}
