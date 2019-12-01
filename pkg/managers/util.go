package managers


func GetManagerName(crName string, deviceName string) string {
	return crName + "-" + deviceName + "-manager"
}
