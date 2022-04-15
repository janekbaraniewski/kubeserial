package api

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("ApiClient")

type ApiClient struct {
	Client client.Client
	Scheme *runtime.Scheme
}

func (c *ApiClient) EnsureConfigMap(ctx context.Context, cr *appv1alpha1.KubeSerial, cm *corev1.ConfigMap) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)

	if err := controllerutil.SetControllerReference(cr, cm, c.Scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}

	found := &corev1.ConfigMap{}
	err := c.Client.Get(ctx, types.NamespacedName{Name: cm.Name, Namespace: cm.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new ConfigMap " + cm.Name)
		err = c.Client.Create(ctx, cm)
		if err != nil {
			logger.Info("ConfigMap not created")
			return err
		}
	} else if err != nil {
		logger.Info("ConfigMap not found")
		return err
	}

	return nil
}

func (r *ApiClient) EnsureService(ctx context.Context, cr *appv1alpha1.KubeSerial, svc *corev1.Service) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)

	if err := controllerutil.SetControllerReference(cr, svc, r.Scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}

	found := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Service" + svc.Name)
		err = r.Client.Create(ctx, svc)
		if err != nil {
			logger.Info("Service not created")
			return err
		}
	} else if err != nil {
		logger.Info("Service not found")
		return err
	}

	return nil
}

func (r *ApiClient) EnsureIngress(ctx context.Context, cr *appv1alpha1.KubeSerial, ingress *networkingv1.Ingress) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)

	if err := controllerutil.SetControllerReference(cr, ingress, r.Scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}

	found := &networkingv1.Ingress{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Ingress " + ingress.Name)
		err = r.Client.Create(ctx, ingress)
		if err != nil {
			logger.Info("Deployment not created")
			return err
		}
	} else if err != nil {
		logger.Info("Deployment not found")
		return err
	}

	return nil
}

func (r *ApiClient) EnsureDeployment(ctx context.Context, cr *appv1alpha1.KubeSerial, deployment *appsv1.Deployment) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)

	if err := controllerutil.SetControllerReference(cr, deployment, r.Scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}

	found := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new Deployment " + deployment.Name)
		err = r.Client.Create(ctx, deployment)
		if err != nil {
			logger.Info("Deployment not created")
			return err
		}
	} else if err != nil {
		logger.Info("Deployment not found")
		return err
	}

	return nil
}

func (r *ApiClient) EnsureDaemonSet(ctx context.Context, cr *appv1alpha1.KubeSerial, ds *appsv1.DaemonSet) error {
	if err := controllerutil.SetControllerReference(cr, ds, r.Scheme); err != nil {
		return err
	}

	found := &appsv1.DaemonSet{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.Client.Create(ctx, ds)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (r *ApiClient) DeleteDeployment(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	d := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.Namespace}, d)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, d, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}

func (r *ApiClient) DeleteConfigMap(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.Namespace}, cm)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, cm, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}

func (r *ApiClient) DeleteService(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	svc := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.Namespace}, svc)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, svc, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil

}

func (r *ApiClient) DeleteIngress(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	ingress := &networkingv1.Ingress{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.Namespace}, ingress)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, ingress, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}
