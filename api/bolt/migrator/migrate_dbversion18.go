package migrator

import portainer "github.com/cloudogu/portainer-ce/api"

func (m *Migrator) updateSettingsToDBVersion19() error {
	legacySettings, err := m.settingsService.Settings()
	if err != nil {
		return err
	}

	if legacySettings.EdgeAgentCheckinInterval == 0 {
		legacySettings.EdgeAgentCheckinInterval = portainer.DefaultEdgeAgentCheckinIntervalInSeconds
	}

	return m.settingsService.UpdateSettings(legacySettings)
}
