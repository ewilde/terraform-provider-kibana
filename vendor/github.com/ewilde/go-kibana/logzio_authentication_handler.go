package kibana

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/xlzd/gotp"
)

var mfaCodeExpiredError = errors.New("the mfa code sent is expired")

func NewLogzAuthenticationHandler(agent *gorequest.SuperAgent) *LogzAuthenticationHandler {
	return &LogzAuthenticationHandler{}
}

func (auth *LogzAuthenticationHandler) Initialize(agent *gorequest.SuperAgent) error {
	if auth.sessionToken != "" {
		auth.setLogzHeaders(agent)
		return nil
	}

	if auth.MfaSecret != "" {
		return auth.initializeWithMFA(agent)
	} else {
		return auth.initializeWithAuth0(agent)
	}
}

func (auth *LogzAuthenticationHandler) initializeWithAuth0(agent *gorequest.SuperAgent) error {
	csrfToken, err := auth.getCSRFToken()

	if err != nil {
		return err
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

func (auth *LogzAuthenticationHandler) initializeWithMFA(agent *gorequest.SuperAgent) error {
	request := gorequest.New()
	csrfToken, err := auth.getCSRFToken()

	if err != nil {
		return err
	}

	mfaCode, secondsLeftForMfaToExpire := auth.getMfaCodeWithExpiry()
	sessionToken, err := auth.getLogzioSessionToken(mfaCode)

	// Attempt regeneration if possible
	if err == mfaCodeExpiredError && secondsLeftForMfaToExpire < 5 {
		log.Print("The mfa code was too close to expiry, so we re-generate and try again")
		mfaCode, secondsLeftForMfaToExpire = auth.getMfaCodeWithExpiry()
		sessionToken, err = auth.getLogzioSessionToken(mfaCode)
	}

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

func (auth *LogzAuthenticationHandler) getLogzioSessionToken(mfaCode string) (string, error) {

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
	  "client_id": "%s",
	  "mfa_code": "%s"
	}`, auth.UserName, auth.Password, auth.ClientId, mfaCode)).
		End()

	if errs != nil {
		return "", errs[0]
	}

	if response.StatusCode == 401 && mfaCode != "" {
		return "", mfaCodeExpiredError
	}

	if response.StatusCode >= 300 {
		return "", errors.New(fmt.Sprintf("Status: %d, %s", response.StatusCode, body))
	}

	authResponse := &Auth0Response{}
	if err := json.Unmarshal([]byte(body), authResponse); err != nil {
		return "", err
	}

	return authResponse.IdTokens, nil

}
func (auth *LogzAuthenticationHandler) getMfaCodeWithExpiry() (string, int64) {
	mfaCode, otpExpirationTime := gotp.NewDefaultTOTP(auth.MfaSecret).NowWithExpiration()
	secondsForOtpExpiry := otpExpirationTime - time.Now().Unix()
	return mfaCode, secondsForOtpExpiry
}

func (auth *LogzAuthenticationHandler) setLogzHeaders(agent *gorequest.SuperAgent) *gorequest.SuperAgent {
	return agent.
		Set("x-logz-csrf-token", auth.csrfToken).
		Set("cookie", fmt.Sprintf("Logzio-Csrf=%s", auth.csrfToken)).
		Set("x-auth-token", auth.sessionToken).
		Set("Content-Type", "application/json")
}
