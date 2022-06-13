package kubeserial

type ResourceSpecPath string
type ResourceLabel string

// Resource config paths
const (
	GatewayCMSpecPath     ResourceSpecPath = "/config/gateway-configmap.yaml"
	GatewayDeploySpecPath ResourceSpecPath = "/config/gateway-deployment.yaml"
	GatewaySvcSpecPath    ResourceSpecPath = "/config/gateway-service.yaml"

	ManagerCMSpecPath     ResourceSpecPath = "/config/manager-configmap.yaml"
	ManagerDeploySpecPath ResourceSpecPath = "/config/manager-deployment.yaml"
	ManagerSvcSpecPath    ResourceSpecPath = "/config/manager-service.yaml"

	MonitorCMSpecPath ResourceSpecPath = "/config/monitor-configmap.yaml"
	MonitorDSSpecPath ResourceSpecPath = "/config/monitor-daemonset.yaml"
)

const (
	AppNameLabel ResourceLabel = "app.kubernetes.io/name"
)
