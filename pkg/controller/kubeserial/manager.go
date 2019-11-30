package kubeserial

import (
	erro "errors"  // TODO: fix this alias

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
)

func (r *ReconcileKubeSerial) ReconcileManager(cr *appv1alpha1.KubeSerial) error {
	for _, device := range cr.Spec.Devices {
		stateCM, err := r.GetDeviceState(&device, cr)
		if err != nil {
			return err
		}

		if stateCM.Data["available"] == "true" {
			if err := r.scheduleManager(cr, &device); err != nil {
				return err
			}
		} else {
			if err := r.deleteManager(cr, &device); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *ReconcileKubeSerial) scheduleManager(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) error {
	switch device.Manager {
	case "octoprint":  // TODO: create some constant
		return r.ScheduleOctoprint(cr, device)
	default:
		return erro.New("Manager not supported")
	}
	return nil
}

func (r *ReconcileKubeSerial) deleteManager(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) error {
	switch device.Manager {
	case "octoprint":  // TODO: create some constant
		return r.DeleteOctoprint(cr, device)
	default:
		return erro.New("Manager not supported") // TODO: add support for custom definitions
	}
	return nil
}
