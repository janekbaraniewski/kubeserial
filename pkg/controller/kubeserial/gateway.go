package kubeserial

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/gateway"
	"github.com/janekbaraniewski/kubeserial/pkg/controller/api"
)



func (r *ReconcileKubeSerial) ReconcileGateway(cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	logger := log.WithValues("Namespace", cr.Namespace, "Name", cr.Name)

	for _, device := range cr.Spec.Devices {
		logger.Info("Reconciling " + device.Name)

		stateCM, err := r.GetDeviceState(&device, cr)
		if err != nil {
			return err
		}

		if stateCM.Data["available"] == "true" {
			logger.Info(device.Name, "available on node", stateCM.Data["node"])
			if err := r.scheduleGateway(cr, &device, stateCM.Data["node"], api); err != nil {
				return err
			}
		} else {
			logger.Info(device.Name, "unavailable")
			if err := r.deleteGateway(cr, &device, api); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ReconcileKubeSerial) scheduleGateway(cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device, node string, api *api.ApiClient) error {
	cm := gateway.CreateConfigMap(cr, p)
	deploy := gateway.CreateDeployment(cr, p, node)
	svc := gateway.CreateService(cr, p)

	if err := api.EnsureConfigMap(cr, cm); err != nil {
		return err
	}

	if err := api.EnsureDeployment(cr, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(cr, svc); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileKubeSerial) deleteGateway(cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device, api *api.ApiClient) error {
	logger := log.WithValues("Component", "gateway", "func", "ensureNotRunning")

	logger.Info("Device " + p.Name)

	name := cr.Name + "-" + p.Name + "-gateway"  // TODO: move name generation to some utils so it's in one place

	if err := api.DeleteDeployment(cr, name); err != nil {
		return err
	}

	if err := api.DeleteConfigMap(cr, name); err != nil {
		return err
	}

	if err := api.DeleteService(cr, name); err != nil {
		return err
	}

	return nil
}
