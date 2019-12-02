package managers


func (m *Manager)GetName(crName string, deviceName string) string {
	return crName + "-" + deviceName + "-manager"
}
