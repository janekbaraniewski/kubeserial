package controllers

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
)

func (r *KubeSerialReconciler) ReconcileMonitor(cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	conf := monitor.CreateConfigMap(cr)
	monitorDaemon := monitor.CreateDaemonSet(cr)

	if err := api.EnsureConfigMap(cr, conf); err != nil {
		return err
	}

	if err := api.EnsureDaemonSet(cr, monitorDaemon); err != nil {
		return err
	}

	return nil
}
