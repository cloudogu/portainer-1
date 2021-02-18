package jwt

import (
	"errors"
	"fmt"
	"github.com/cloudogu/portainer-ce/api"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
)

const (
	blocklistTokenTTLinSeconds = 86400
)

// Service represents a service for managing JWT tokens.
type Service struct {
	secret             []byte
	userSessionTimeout time.Duration
	tokenBlockList     *BlocklistTokenMap
}

type claims struct {
	UserID     int    `json:"id"`
	Username   string `json:"username"`
	Role       int    `json:"role"`
	OAuthToken string `json:"oauth_token"`
	jwt.StandardClaims
}

var (
	errSecretGeneration = errors.New("Unable to generate secret key")
	errInvalidJWTToken  = errors.New("Invalid JWT token")
)

// NewService initializes a new service. It will generate a random key that will be used to sign JWT tokens.
func NewService(userSessionDuration string) (*Service, error) {
	userSessionTimeout, err := time.ParseDuration(userSessionDuration)
	if err != nil {
		return nil, err
	}

	secret := securecookie.GenerateRandomKey(32)
	if secret == nil {
		return nil, errSecretGeneration
	}

	tokenBlockList := NewBlocklistTokenMap(blocklistTokenTTLinSeconds, time.Hour)

	service := &Service{
		secret,
		userSessionTimeout,
		tokenBlockList,
	}
	return service, nil
}

// GenerateToken generates a new JWT token.
func (service *Service) GenerateToken(data *portainer.TokenData) (string, error) {
	expireToken := time.Now().Add(service.userSessionTimeout).Unix()
	cl := claims{
		UserID:     int(data.ID),
		Username:   data.Username,
		Role:       int(data.Role),
		OAuthToken: data.OAuthToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	signedToken, err := token.SignedString(service.secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseAndVerifyToken parses a JWT token and verify its validity. It returns an error if token is invalid.
func (service *Service) ParseAndVerifyToken(token string) (*portainer.TokenData, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, msg
		}
		return service.secret, nil
	})
	if err == nil && parsedToken != nil {
		if cl, ok := parsedToken.Claims.(*claims); ok && parsedToken.Valid {
			tokenData := &portainer.TokenData{
				ID:         portainer.UserID(cl.UserID),
				Username:   cl.Username,
				Role:       portainer.UserRole(cl.Role),
				OAuthToken: cl.OAuthToken,
			}

			if tokenData.OAuthToken != "" {
				if service.tokenBlockList.IsBlocked(tokenData.OAuthToken) {
					return nil, errInvalidJWTToken
				}
			}

			return tokenData, nil
		}
	}
	return nil, errInvalidJWTToken
}

// SetUserSessionDuration sets the user session duration
func (service *Service) SetUserSessionDuration(userSessionDuration time.Duration) {
	service.userSessionTimeout = userSessionDuration
}

// AddTokenToBlocklist adds a token identifier to the token list
func (service *Service) AddTokenToBlocklist(token string) {
	service.tokenBlockList.Put(token)
}
