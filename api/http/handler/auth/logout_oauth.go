package auth

import (
	"fmt"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	portError "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type verifyResponse struct {
	Valid bool `json:"valid"`
}

func (handler *Handler) invalidateOAuthSession(w http.ResponseWriter, r *http.Request) *portError.HandlerError {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &portError.HandlerError{http.StatusInternalServerError, "Unable to block token: ", err}
	}

	content, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return &portError.HandlerError{http.StatusInternalServerError, "Unable to find content-type ", err}
	}

	if content == "application/x-www-form-urlencoded" || content == "text/plain" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return &portError.HandlerError{http.StatusInternalServerError, "Cannot parse body ", err}
		}
		token := values.Get("logoutRequest")

		jwtBlocklist, ok := handler.JWTService.(dataservices.BlocklistedJWTService)
		fmt.Printf("###### token: %s and service could be casted: %t", token, ok)
		if ok {
			jwtBlocklist.AddTokenToBlocklist(token)
		}
		return nil
	} else {
		return &portError.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("Invalid content type %s. Expected \"application/x-www-form-urlencoded\" or \"text/plain\"", content),
		}
	}
}

func (handler *Handler) isJWTTokenNotBlocked(w http.ResponseWriter, r *http.Request) *portError.HandlerError {
	token, handlerErr := handler.retrieveAuthTokenFromRequest(r)
	if handlerErr != nil {
		return handlerErr
	}

	var err error
	_, err = handler.JWTService.ParseAndVerifyToken(token)
	if err != nil {
		return response.JSON(w, &verifyResponse{Valid: false})
	}

	return response.JSON(w, &verifyResponse{Valid: true})
}

func (handler *Handler) retrieveAuthTokenFromRequest(r *http.Request) (string, *portError.HandlerError) {
	// Optionally, token might be set via the "token" query parameter.
	// For example, in websocket requests
	token := r.URL.Query().Get("token")

	// Get token from the Authorization header
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		return "", &portError.HandlerError{StatusCode: http.StatusBadRequest, Message: "Missing token"}
	}
	return token, nil
}
