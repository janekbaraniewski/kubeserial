package managers

import "strings"

func (m *Manager) GetName(crName string, deviceName string) string {
	return strings.ToLower(crName + "-" + deviceName + "-manager")
}
