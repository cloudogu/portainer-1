package migrator

import (
	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/internal/authorization"

	"github.com/rs/zerolog/log"
)

func (m *Migrator) migrateDBVersionToDB36() error {
	log.Info().Msg("updating user authorizations")

	return m.migrateUsersToDB36()
}

func (m *Migrator) migrateUsersToDB36() error {
	log.Info().Msg("updating user authorizations")

	users, err := m.userService.Users()
	if err != nil {
		return err
	}

	for _, user := range users {
		currentAuthorizations := authorization.DefaultPortainerAuthorizations()
		currentAuthorizations[portainer.OperationPortainerUserListToken] = true
		currentAuthorizations[portainer.OperationPortainerUserCreateToken] = true
		currentAuthorizations[portainer.OperationPortainerUserRevokeToken] = true
		user.PortainerAuthorizations = currentAuthorizations
		err = m.userService.UpdateUser(user.ID, &user)
		if err != nil {
			return err
		}
	}

	return nil
}
