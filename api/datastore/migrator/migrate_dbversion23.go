package migrator

import (
	portainer "github.com/cloudogu/portainer-ce/api"

	"github.com/rs/zerolog/log"
)

func (m *Migrator) updateSettingsToDB24() error {
	log.Info().Msg("updating Settings")

	legacySettings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	legacySettings.AllowHostNamespaceForRegularUsers = true
	legacySettings.AllowDeviceMappingForRegularUsers = true
	legacySettings.AllowStackManagementForRegularUsers = true

	return m.settingsService.UpdateSettings(legacySettings)
}

func (m *Migrator) updateStacksToDB24() error {
	log.Info().Msg("updating stacks")

	stacks, err := m.stackService.Stacks()
	if err != nil {
		return err
	}

	for idx := range stacks {
		stack := &stacks[idx]
		stack.Status = portainer.StackStatusActive
		err := m.stackService.UpdateStack(stack.ID, stack)
		if err != nil {
			return err
		}
	}

	return nil
}
