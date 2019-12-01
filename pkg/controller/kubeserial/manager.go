package kubeserial

import (
	"errors"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controller/api"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
	"github.com/janekbaraniewski/kubeserial/pkg/managers/octoprint"
)

func (r *ReconcileKubeSerial) ReconcileManagers(cr *appv1alpha1.KubeSerial, api *api.ApiClient) error {
	for _, device := range cr.Spec.Devices {
		stateCM, err := r.GetDeviceState(&device, cr)
		if err != nil {
			return err
		}

		if stateCM.Data["available"] == "true" {
			if err := r.scheduleManager(cr, &device, api); err != nil {
				return err
			}
		} else {
			if err := r.deleteManager(cr, &device, api); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReconcileKubeSerial) scheduleManager(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	switch device.Manager {
	case "octoprint":  // TODO: create some constant
		if err := octoprint.Schedule(cr, device, api); err != nil {
			return err
		}
	default:
		return errors.New("Manager not supported")
	}

	if cr.Spec.Ingress.Enabled {
		ingress := managers.CreateIngress(cr, device, cr.Spec.Ingress.Domain)
		if err := api.EnsureIngress(cr, ingress); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileKubeSerial) deleteManager(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	switch device.Manager {
	case "octoprint":  // TODO: create some constant
		return octoprint.Delete(cr, device, api)
	default:
		return errors.New("Manager not supported") // TODO: add support for custom definitions
	}


	name := cr.Name + "-" + device.Name + "-manager"

	if err := api.DeleteIngress(cr, name); err != nil {
		return err
	}

	return nil
}
