package controllers

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
)

func (r *KubeSerialReconciler) ReconcileManagers(ctx context.Context, cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	for _, device := range cr.Spec.Devices {
		stateCM, err := r.GetDeviceState(ctx, &device, cr)
		if err != nil {
			return err
		}
		manager := managers.Available[device.Manager]
		if stateCM.Data["available"] == "true" {
			if err := manager.Schedule(ctx, cr, &device, api); err != nil {
				return err
			}
		} else {
			if err := manager.Delete(ctx, cr, &device, api); err != nil {
				return err
			}
		}
	}
	return nil
}
