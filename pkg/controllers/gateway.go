package controllers

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/gateway"
)

func (r *KubeSerialReconciler) ReconcileGateway(ctx context.Context, cr *appv1alpha1.KubeSerial, api api.API) error {
	logger := log.WithName("ReconcileGateway")
	for _, device := range cr.Spec.Devices {
		logger.Info("Reconciling " + device.Name)

		stateCM, err := r.GetDeviceState(ctx, &device, cr)
		if err != nil {
			return err
		}

		if stateCM.Data["available"] == "true" {
			logger.Info("available on node", "NodeName", stateCM.Data["node"], "DeviceName", device.Name)
			if err := r.scheduleGateway(ctx, cr, &device, stateCM.Data["node"], api); err != nil {
				return err
			}
		} else {
			logger.Info("unavailable", "DeviceName", device.Name)
			if err := r.deleteGateway(ctx, cr, &device, api); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *KubeSerialReconciler) scheduleGateway(ctx context.Context, cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device, node string, api api.API) error {
	cm := gateway.CreateConfigMap(cr, p)
	deploy := gateway.CreateDeployment(cr, p, node)
	svc := gateway.CreateService(cr, p)

	if err := api.EnsureConfigMap(ctx, cr, cm); err != nil {
		return err
	}

	if err := api.EnsureDeployment(ctx, cr, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(ctx, cr, svc); err != nil {
		return err
	}

	return nil
}

func (r *KubeSerialReconciler) deleteGateway(ctx context.Context, cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device, api api.API) error {
	logger := log.WithName("deleteGateway")

	logger.Info("Device " + p.Name)

	name := cr.Name + "-" + p.Name + "-gateway" // TODO: move name generation to some utils so it's in one place

	if err := api.DeleteDeployment(ctx, cr, name); err != nil {
		return err
	}

	if err := api.DeleteConfigMap(ctx, cr, name); err != nil {
		return err
	}

	if err := api.DeleteService(ctx, cr, name); err != nil {
		return err
	}

	return nil
}
