package octoprint

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) *corev1.ConfigMap {
	labels := map[string]string{
		"app": cr.Name + "-" + device.Name + "-manager",
	}
  // TODO: parametrize
	conf := `
accessControl:
  enabled: false
plugins:
  announcements:
    _config_version: 1
    channels:
      _blog:
        read_until: 1573642500
      _important:
        read_until: 1521111600
      _octopi:
        read_until: 1573722900
      _plugins:
        read_until: 1573862400
      _releases:
        read_until: 1574699400
  discovery:
    upnpUuid: ef35acc7-a859-4947-980d-d5edb10508e4
  softwareupdate:
    _config_version: 6
  tracking:
    enabled: false
deviceProfiles:
  default: _default
serial:
  additionalPorts:
  - /dev/device
  autoconnect: true
  baudrate: 0
  port: /dev/device
server:
  firstRun: false
  onlineCheck:
    enabled: true
  pluginBlacklist:
    enabled: false
  seenWizards:
    corewizard: 3
    cura: null
    tracking: null`

	return &corev1.ConfigMap {
		ObjectMeta:		metav1.ObjectMeta {
			Name: 		cr.Name + "-" + device.Name + "-manager",
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Data:			map[string]string {
			"config.yaml": conf,
		},
	}
}
