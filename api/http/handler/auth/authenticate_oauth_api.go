package auth

import (
	"errors"
	"github.com/portainer/libhttp/request"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/cloudogu/portainer-ce/api"
	httperror "github.com/portainer/libhttp/error"
)

type authenticateApiPayload struct {
	Username string
	Password string
}

type authenticateApiResponse struct {
	JWT string `json:"jwt"`
}

func (payload *authenticateApiPayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Username) {
		return errors.New("invalid username")
	}
	if govalidator.IsNull(payload.Password) {
		return errors.New("invalid password")
	}
	return nil
}

func (handler *Handler) authenticateViaApi(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload authenticatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return &httperror.HandlerError{http.StatusBadRequest, "Invalid request payload", err}
	}

	settings, err := handler.DataStore.Settings().Settings()
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to retrieve settings from the database", Err: err}
	}

	user, err := handler.DataStore.User().UserByUsername(payload.Username)
	if err != nil && !handler.DataStore.IsErrObjectNotFound(err) {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to retrieve a user with the specified username from the database", Err: err}
	}

	if settings.AuthenticationMethod != portainer.AuthenticationOAuth {
		return &httperror.HandlerError{StatusCode: http.StatusUnprocessableEntity, Message: "Invalid authentication method"}
	}

	//1) Request ticket with user login
	u, _ := url.ParseRequestURI(settings.OAuthSettings.LogoutURI)
	u.Path = "/cas/v1/tickets"
	urlStr := u.String()
	log.Print("TicketURL: " + urlStr)

	client := &http.Client{}
	data := url.Values{}
	data.Set("username", payload.Username)
	data.Set("password", payload.Password)
	loginRequest, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to create cas login request", Err: err}
	}

	loginRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	loginResponse, err := client.Do(loginRequest)
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to send cas login request", Err: err}
	}

	//Use location to create a service ticket
	location, _ := url.ParseRequestURI(loginResponse.Header.Get("Location"))
	data = url.Values{}
	data.Set("service", settings.OAuthSettings.RedirectURI)
	stRequest, err := http.NewRequest("POST", location.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to create cas login request", Err: err}
	}

	stRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	stResponse, err := client.Do(stRequest)
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to send cas st ticket request", Err: err}
	}

	defer r.Body.Close()
	body, err := io.ReadAll(stResponse.Body)
	if err != nil {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Unable to block token: ", Err: err}
	}

	if !strings.Contains(string(body), "ST-") {
		return &httperror.HandlerError{StatusCode: http.StatusInternalServerError, Message: "Invalid service ticket returned: ", Err: err}
	}

	user.OAuthToken = &oauth2.Token{AccessToken: string(body)}

	return handler.writeToken(w, user, false)
}
