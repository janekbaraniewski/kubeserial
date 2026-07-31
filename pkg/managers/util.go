package managers

import "github.com/janekbaraniewski/kubeserial/pkg/utils/apis"

func (m *Manager) GetName(crName string, deviceName string) string {
	return apis.ManagerName(crName, deviceName)
}
