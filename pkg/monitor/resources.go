package monitor

import (
	"fmt"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func CreateConfigMap(fs utils.FileSystem, devices []appv1alpha1.SerialDevice2) (*corev1.ConfigMap, error) {
	rule := ""
	for _, device := range devices {
		rule += fmt.Sprintf(
			"SUBSYSTEM==\"tty\", ATTRS{IDVendor}==\"%s\", ATTRS{IDProduct}==\"%s\", SYMLINK+=\"%s\"\n",
			device.IDVendor,
			device.IDProduct,
			device.Name,
		)
	}

	cm := &corev1.ConfigMap{}

	if err := utils.LoadResourceFromYaml(fs, kubeserial.MonitorCMSpecPath, cm); err != nil {
		return cm, err
	}

	cm.Data["98-devices.rules"] = rule

	return cm, nil
}

func CreateDaemonSet(fs utils.FileSystem) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	if err := utils.LoadResourceFromYaml(fs, kubeserial.MonitorDSSpecPath, ds); err != nil {
		return ds, err
	}

	return ds, nil
}
