package monitor

import (
	"fmt"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func CreateConfigMap(fs utils.FileSystem, devices []appv1alpha1.SerialDevice_2) (*corev1.ConfigMap, error) {
	SPEC_PATH := "/config/monitor-configmap.yaml"

	rule := ""
	for _, device := range devices {
		rule += fmt.Sprintf(
			"SUBSYSTEM==\"tty\", ATTRS{idVendor}==\"%s\", ATTRS{idProduct}==\"%s\", SYMLINK+=\"%s\"\n",
			device.IdVendor,
			device.IdProduct,
			device.Name,
		)
	}

	cm := &corev1.ConfigMap{}

	if err := utils.LoadResourceFromYaml(fs, SPEC_PATH, cm); err != nil {
		return cm, err
	}

	cm.Data["98-devices.rules"] = rule

	return cm, nil
}

func CreateDaemonSet(fs utils.FileSystem) (*appsv1.DaemonSet, error) {
	SPEC_PATH := "/config/monitor-daemonset.yaml"

	ds := &appsv1.DaemonSet{}

	if err := utils.LoadResourceFromYaml(fs, SPEC_PATH, ds); err != nil {
		return ds, err
	}

	return ds, nil
}