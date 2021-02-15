package oauth

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"

	"github.com/portainer/portainer/api"
)

// Service represents a service used to authenticate users against an authorization server
type Service struct{}

// NewService returns a pointer to a new instance of this service
func NewService() *Service {
	return &Service{}
}

type authenticationData struct {
	ID         string `json:"id"`
	Attributes struct {
		Username    string   `json:"username"`
		Cn          string   `json:"cn"`
		Mail        string   `json:"mail"`
		GivenName   string   `json:"givenName"`
		Surname     string   `json:"surname"`
		DisplayName string   `json:"displayName"`
		Groups      []string `json:"groups"`
	} `json:"attributes"`
}

// Authenticate takes an access code and exchanges it for an access token from portainer OAuthSettings token endpoint.
// On success, it will then return the username associated to authenticated user by fetching this information
// from the resource server and matching it with the user identifier setting.
func (*Service) Authenticate(code string, configuration *portainer.OAuthSettings) (portainer.OAuthUserData, error) {
	token, err := getAccessToken(code, configuration)
	if err != nil {
		log.Printf("[DEBUG] - Failed retrieving access token: %v", err)
		return portainer.OAuthUserData{}, err
	}

	userData, err := getUserData(token, configuration)
	if err != nil {
		log.Printf("[DEBUG] - Failed retrieving username: %v", err)
		return portainer.OAuthUserData{}, err
	}

	return userData, nil
}

func getAccessToken(code string, configuration *portainer.OAuthSettings) (string, error) {
	unescapedCode, err := url.QueryUnescape(code)
	if err != nil {
		return "", err
	}

	config := buildConfig(configuration)
	token, err := config.Exchange(context.Background(), unescapedCode)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func getUserData(token string, configuration *portainer.OAuthSettings) (portainer.OAuthUserData, error) {
	req, err := http.NewRequest("GET", configuration.ResourceURI, nil)
	if err != nil {
		return portainer.OAuthUserData{}, err
	}

	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return portainer.OAuthUserData{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return portainer.OAuthUserData{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return portainer.OAuthUserData{}, &oauth2.RetrieveError{
			Response: resp,
			Body:     body,
		}
	}

	content, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return portainer.OAuthUserData{}, err
	}

	if content == "application/x-www-form-urlencoded" || content == "text/plain" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return portainer.OAuthUserData{}, err
		}

		username := values.Get(configuration.UserIdentifier)
		if username == "" {
			return portainer.OAuthUserData{}, &oauth2.RetrieveError{
				Response: resp,
				Body:     body,
			}
		}

		userData := portainer.OAuthUserData{
			Username:   username,
			OAuthToken: token,
		}
		return userData, nil
	}

	var data authenticationData
	if err = json.Unmarshal(body, &data); err != nil {
		return portainer.OAuthUserData{}, err
	}

	if data.ID != "" {
		userData := portainer.OAuthUserData{
			Username:   data.ID,
			OAuthToken: token,
			Teams:      data.Attributes.Groups,
		}
		return userData, nil
	}

	return portainer.OAuthUserData{}, &oauth2.RetrieveError{
		Response: resp,
		Body:     body,
	}
}

func buildConfig(configuration *portainer.OAuthSettings) *oauth2.Config {
	endpoint := oauth2.Endpoint{
		AuthURL:  configuration.AuthorizationURI,
		TokenURL: configuration.AccessTokenURI,
	}

	return &oauth2.Config{
		ClientID:     configuration.ClientID,
		ClientSecret: configuration.ClientSecret,
		Endpoint:     endpoint,
		RedirectURL:  configuration.RedirectURI,
		Scopes:       []string{configuration.Scopes},
	}
}
