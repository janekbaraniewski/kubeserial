package controllers

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
)

func (r *KubeSerialReconciler) ReconcileManagers(cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	for _, device := range cr.Spec.Devices {
		stateCM, err := r.GetDeviceState(&device, cr)
		if err != nil {
			return err
		}
		manager := managers.Available[device.Manager]
		if stateCM.Data["available"] == "true" {
			if err := manager.Schedule(cr, &device, api); err != nil {
				return err
			}
		} else {
			if err := manager.Delete(cr, &device, api); err != nil {
				return err
			}
		}
	}
	return nil
}
