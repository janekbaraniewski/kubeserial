// Package apis holds helpers for deriving the names of the Kubernetes objects
// kubeserial creates, so the naming scheme lives in one place.
package apis

import "strings"

const (
	gatewaySuffix = "-gateway"
	managerSuffix = "-manager"
)

// ScheduleRequestName returns the name of the ManagerScheduleRequest created
// for the given device/manager pair.
func ScheduleRequestName(deviceName, managerName string) string {
	return deviceName + "-" + managerName
}

// GatewayName returns the name of the gateway objects (config map, deployment,
// service) created for the given device.
func GatewayName(deviceName string) string {
	return deviceName + gatewaySuffix
}

// ManagerName returns the name of the manager objects scheduled for the given
// schedule-request/device pair.
func ManagerName(requestName, deviceName string) string {
	return strings.ToLower(requestName + "-" + deviceName + managerSuffix)
}
