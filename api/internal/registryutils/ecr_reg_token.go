package registryutils

import (
	"time"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/aws/ecr"
	"github.com/cloudogu/portainer-ce/api/dataservices"

	"github.com/rs/zerolog/log"
)

func isRegTokenValid(registry *portainer.Registry) (valid bool) {
	return registry.AccessToken != "" && registry.AccessTokenExpiry > time.Now().Unix()
}

func doGetRegToken(dataStore dataservices.DataStore, registry *portainer.Registry) (err error) {
	ecrClient := ecr.NewService(registry.Username, registry.Password, registry.Ecr.Region)
	accessToken, expiryAt, err := ecrClient.GetAuthorizationToken()
	if err != nil {
		return
	}

	registry.AccessToken = *accessToken
	registry.AccessTokenExpiry = expiryAt.Unix()

	err = dataStore.Registry().UpdateRegistry(registry.ID, registry)

	return
}

func parseRegToken(registry *portainer.Registry) (username, password string, err error) {
	ecrClient := ecr.NewService(registry.Username, registry.Password, registry.Ecr.Region)
	return ecrClient.ParseAuthorizationToken(registry.AccessToken)
}

func EnsureRegTokenValid(dataStore dataservices.DataStore, registry *portainer.Registry) (err error) {
	if registry.Type == portainer.EcrRegistry {
		if isRegTokenValid(registry) {
			log.Debug().Msg("current ECR token is still valid")
		} else {
			err = doGetRegToken(dataStore, registry)
			if err != nil {
				log.Debug().Msg("refresh ECR token")
			}
		}
	}

	return
}

func GetRegEffectiveCredential(registry *portainer.Registry) (username, password string, err error) {
	if registry.Type == portainer.EcrRegistry {
		username, password, err = parseRegToken(registry)
	} else {
		username = registry.Username
		password = registry.Password
	}
	return
}
