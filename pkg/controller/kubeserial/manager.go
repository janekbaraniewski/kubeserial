package kubeserial

import (
	"context"
	erro "errors"  // TODO: fix this alias

	v1beta1 "k8s.io/api/extensions/v1beta1"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	objects "github.com/janekbaraniewski/kubeserial/pkg/managers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		if err := r.ScheduleOctoprint(cr, device); err != nil {
			return err
		}
	default:
		return erro.New("Manager not supported")
	}

	if cr.Spec.Ingress.Enabled {
		ingress := objects.CreateIngress(cr, device, cr.Spec.Ingress.Domain)
		if err := r.ReconcileIngress(cr, ingress); err != nil {
			return err
		}
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


	name := objects.GetManagerName(cr.Name, device.Name)
	ingress := &v1beta1.Ingress{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, ingress)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.client.Delete(context.TODO(), ingress, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}

	return nil
}
