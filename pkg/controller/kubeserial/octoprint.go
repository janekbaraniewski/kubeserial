package kubeserial

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	objects "github.com/janekbaraniewski/kubeserial/pkg/managers/octoprint"
)

func (r *ReconcileKubeSerial)ScheduleOctoprint(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) error {
	cm := objects.CreateConfigMap(cr, device)
	deploy := objects.CreateDeployment(cr, device)
	svc := objects.CreateService(cr, device)
	ingress := objects.CreateIngress(cr, device)

	if err := r.ReconcileConfigMap(cr, cm); err != nil {
		return err
	}

	if err := r.ReconcileDeployment(cr, deploy); err != nil {
		return err
	}

	if err := r.ReconcileService(cr, svc); err != nil {
		return err
	}

	if err := r.ReconcileIngress(cr, ingress); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileKubeSerial)DeleteOctoprint(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) error {
	logger := log.WithValues("Component", "manager", "func", "ensureNotRunning")

	logger.Info("Device " + device.Name)

	name := cr.Name + "-" + device.Name + "-manager"  // TODO: move name generation to some utils so it's in one place

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

	ingress := &v1beta1.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: cr.Namespace}, ingress)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.client.Delete(context.TODO(), ingress, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}


	return nil
}
