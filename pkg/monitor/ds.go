package monitor

import (
	appsv1 "k8s.io/api/apps/v1"

	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

func CreateDaemonSet(fs utils.FileSystem) (*appsv1.DaemonSet, error) {
	SPEC_PATH := "/config/monitor-daemonset-spec.yaml"

	ds := &appsv1.DaemonSet{}

	if err := utils.LoadResourceFromYaml(fs, SPEC_PATH, ds); err != nil {
		return ds, err
	}

	return ds, nil
}
