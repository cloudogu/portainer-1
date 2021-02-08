package auth

import (
	"errors"
	. "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

type verifyResponse struct {
	Valid bool `json:"valid"`
}

func (handler *Handler) invalidateOAuthSession(w http.ResponseWriter, r *http.Request) *HandlerError {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &HandlerError{http.StatusInternalServerError, "Unable to block token: ", err}
	}

	content, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return &HandlerError{http.StatusInternalServerError, "Unable to find content-type ", err}
	}

	if content == "application/x-www-form-urlencoded" || content == "text/plain" {
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return &HandlerError{http.StatusInternalServerError, "Cannot parse body ", err}
		}

		token := values.Get("logoutRequest")

		handler.JWTService.AddTokenToBlacklist(token)

		w.WriteHeader(200)
		return nil
	} else {
		return &HandlerError{http.StatusInternalServerError, "Invalid content type ", err}
	}
}

func (handler *Handler) isJWTTokenNotBlocked(w http.ResponseWriter, r *http.Request) *HandlerError {
	var token string

	// Optionally, token might be set via the "token" query parameter.
	// For example, in websocket requests
	token = r.URL.Query().Get("token")

	// Get token from the Authorization header
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		return &HandlerError{http.StatusBadRequest, "Missing token", errors.New("Missing token")}
	}

	var err error
	_, err = handler.JWTService.ParseAndVerifyToken(token)

	if err != nil {
		return response.JSON(w, &verifyResponse{Valid: false})
	}

	return response.JSON(w, &verifyResponse{Valid: true})
}
