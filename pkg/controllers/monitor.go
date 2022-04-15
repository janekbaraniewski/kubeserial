package controllers

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
)

func (r *KubeSerialReconciler) ReconcileMonitor(ctx context.Context, cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	conf := monitor.CreateConfigMap(cr)
	monitorDaemon := monitor.CreateDaemonSet(cr)

	if err := api.EnsureConfigMap(ctx, cr, conf); err != nil {
		return err
	}

	if err := api.EnsureDaemonSet(ctx, cr, monitorDaemon); err != nil {
		return err
	}

	return nil
}
