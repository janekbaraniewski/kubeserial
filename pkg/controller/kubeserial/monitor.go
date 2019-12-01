package kubeserial

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
	"github.com/janekbaraniewski/kubeserial/pkg/controller/api"
)

func (r *ReconcileKubeSerial) ReconcileMonitor(cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	conf 	:= monitor.CreateConfigMap(cr)
	monitorDaemon 	:= monitor.CreateDaemonSet(cr)

	if err := api.EnsureConfigMap(cr, conf); err != nil {
		return err
	}

	if err := api.EnsureDaemonSet(cr, monitorDaemon); err != nil {
		return err
	}

	return nil
}
