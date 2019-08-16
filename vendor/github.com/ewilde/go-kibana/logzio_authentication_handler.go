package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/xlzd/gotp"
	"log"
	"regexp"
)

const (
	auth0MFAInvalidCode = "a0.mfa_invalid_code"
)

func NewLogzAuthenticationHandler(agent *gorequest.SuperAgent) *LogzAuthenticationHandler {
	return &LogzAuthenticationHandler{}
}

func (auth *LogzAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	if auth.sessionToken != "" {
		auth.setLogzHeaders(agent)
		return nil
	}

	if auth.MfaSecret != "" {
		return auth.initializeWithAuth0MFA(agent)
	} else {
		return auth.initializeWithAuth0(agent)
	}
}

// auth0RO sends an Auth0 resource owner request with the given form
func (auth *LogzAuthenticationHandler) auth0RO(form string) (response *Auth0Response, err error) {
	request := gorequest.New()
	rawResponse, body, errs := request.Post(fmt.Sprintf("%s/oauth/ro", auth.Auth0Uri)).
		Set("kbn-version", DefaultKibanaVersion553).
		Set("Content-Type", "application/x-www-form-urlencoded").
		Type("form").
		Send(form).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	authResponse := &Auth0Response{}
	if err := json.Unmarshal([]byte(body), authResponse); err != nil {
		return nil, err
	}

	// Check for Auth0 errors
	if authResponse.Error != "" {
		error := fmt.Sprintf("Status: %d, Error: %s", rawResponse.StatusCode, authResponse.Error)
		if authResponse.ErrorDescription != "" {
			error += fmt.Sprintf(", Description: %s", authResponse.ErrorDescription)
		}

		error += fmt.Sprintf("\nResponse Body: %s", body)

		// We still return authResponse here in case the caller wants to access the Auth0 errors
		//	e.g. to retry on MFA expiry
		return authResponse, errors.New(error)
	}

	return authResponse, nil
}

// initializeWithAuth0 exchanges non-MFA credentials for a session token
func (auth *LogzAuthenticationHandler) initializeWithAuth0(agent *gorequest.SuperAgent) error {
	csrfToken, err := auth.getCSRFToken()

	if err != nil {
		return err
	}

	form := fmt.Sprintf(`{
	  "scope": "openid email connection",
	  "response_type": "code",
	  "connection": "Username-Password-Authentication",
	  "username": "%s",
	  "password": "%s",
	  "grant_type": "password",
	  "client_id": "%s"
	}`, auth.UserName, auth.Password, auth.ClientId)
	authResponse, err := auth.auth0RO(form)
	if err != nil {
		return err
	}

	request := gorequest.New()
	response, body, errs := request.Post(fmt.Sprintf("%s/login/jwt", auth.LogzUri)).
		Set("x-logz-csrf-token", csrfToken).
		Set("cookie", fmt.Sprintf("Logzio-Csrf=%s", csrfToken)).
		Send(fmt.Sprintf(`{
		  "jwt": "%s"
		}`, authResponse.IdTokens)).
		End()

	if errs != nil {
		return errs[0]
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("error logging in (%d). %s", response.StatusCode, string(body))
	}

	jwtResponse := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &jwtResponse); err != nil {
		return err
	}

	auth.sessionToken = jwtResponse["sessionToken"].(string)
	auth.setLogzHeaders(agent)
	return nil
}

// initializeWithAuth0MFA exchanges MFA credentials for a session token
func (auth *LogzAuthenticationHandler) initializeWithAuth0MFA(agent *gorequest.SuperAgent) error {
	request := gorequest.New()
	csrfToken, err := auth.getCSRFToken()

	if err != nil {
		return err
	}

	sessionToken, err := auth.getLogzioSessionToken(true)
	// If we're still failing, we cannot proceed
	if err != nil {
		return fmt.Errorf("Error getting MFA code: %s", err)
	}

	response, body, errs := request.Post(fmt.Sprintf("%s/login/jwt", auth.LogzUri)).
		Set("x-logz-csrf-token", csrfToken).
		Set("cookie", fmt.Sprintf("Logzio-Csrf=%s", csrfToken)).
		Send(fmt.Sprintf(`{
	  "jwt": "%s"
	}`, sessionToken)).
		End()

	if response.StatusCode >= 400 {
		return fmt.Errorf("error logging in (%d). %s", response.StatusCode, string(body))
	}

	if errs != nil {
		return errs[0]
	}

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

func (auth *LogzAuthenticationHandler) getCSRFToken() (string, error) {
	request := gorequest.New()
	response, _, errs := request.Get(fmt.Sprintf("%s/#/login", auth.LogzUri)).
		End()

	if len(errs) > 0 {
		return "", errs[0]
	}
	cookieHeader := response.Header.Get("Set-Cookie")
	csrfCookieRegEx := regexp.MustCompile("Logzio-Csrf=([^;]+)")
	cookieRegExMatches := csrfCookieRegEx.FindStringSubmatch(cookieHeader)
	if len(cookieRegExMatches) < 2 {
		return "", errors.New("could not retrieve CSRF token from logz.io cookie")
	}
	csrfToken := cookieRegExMatches[1]
	auth.csrfToken = csrfToken
	return csrfToken, nil
}

func (auth *LogzAuthenticationHandler) getLogzioSessionToken(retry bool) (sessionToken string, err error) {
	mfaCode := auth.getMFACode()

	form := fmt.Sprintf(`{
	  "scope": "openid email connection",
	  "response_type": "code",
	  "connection": "Username-Password-Authentication",
	  "username": "%s",
	  "password": "%s",
	  "grant_type": "password",
	  "client_id": "%s",
	  "mfa_code": "%s"
	}`, auth.UserName, auth.Password, auth.ClientId, mfaCode)
	authResponse, err := auth.auth0RO(form)

	if authResponse != nil && authResponse.Error == auth0MFAInvalidCode && retry {
		log.Print("MFA code potentially expired, so we re-generate and try again")
		sessionToken, err = auth.getLogzioSessionToken(false)

		if err != nil {
			return
		}
	} else if err != nil {
		return
	}

	sessionToken = authResponse.IdTokens
	return
}

func (auth *LogzAuthenticationHandler) getMFACode() string {
	return gotp.NewDefaultTOTP(auth.MfaSecret).Now()
}

func (auth *LogzAuthenticationHandler) setLogzHeaders(agent *gorequest.SuperAgent) *gorequest.SuperAgent {
	return agent.
		Set("x-logz-csrf-token", auth.csrfToken).
		Set("cookie", fmt.Sprintf("Logzio-Csrf=%s", auth.csrfToken)).
		Set("x-auth-token", auth.sessionToken).
		Set("Content-Type", "application/json")
}
