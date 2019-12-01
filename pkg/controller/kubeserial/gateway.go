package kubeserial

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)



func (r *ReconcileKubeSerial) ReconcileGateway(cr *appv1alpha1.KubeSerial) error {
	logger := log.WithValues("Namespace", cr.Namespace, "Name", cr.Name)

	for _, device := range cr.Spec.Devices {
		logger.Info("Reconciling " + device.Name)

		stateCM, err := r.GetDeviceState(&device, cr)
		if err != nil {
			return err
		}

		if stateCM.Data["available"] == "true" {
			logger.Info(device.Name, "available on node", stateCM.Data["node"])
			if err := r.scheduleGateway(cr, &device, stateCM.Data["node"]); err != nil {
				return err
			}
		} else {
			logger.Info(device.Name, "unavailable")
			if err := r.deleteGateway(cr, &device); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ReconcileKubeSerial) scheduleGateway(cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device, node string) error {
	cm := gateway.CreateConfigMap(cr, p)
	deploy := gateway.CreateDeployment(cr, p, node)
	svc := gateway.CreateService(cr, p)

	if err := r.ReconcileConfigMap(cr, cm); err != nil {
		return err
	}

	if err := r.ReconcileDeployment(cr, deploy); err != nil {
		return err
	}

	if err := r.ReconcileService(cr, svc); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileKubeSerial) deleteGateway(cr *appv1alpha1.KubeSerial, p *appv1alpha1.Device) error {
	logger := log.WithValues("Component", "gateway", "func", "ensureNotRunning")

	logger.Info("Device " + p.Name)

	name := cr.Name + "-" + p.Name + "-gateway"  // TODO: move name generation to some utils so it's in one place

	d := &v1beta2.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, d)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.client.Delete(context.TODO(), d, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}

	cm := &corev1.ConfigMap{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, cm)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.client.Delete(context.TODO(), cm, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}

	s := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, s)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.client.Delete(context.TODO(), s, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}

	return nil
}
