package kubeapi

import (
	"context"

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

type API interface {
	EnsureObject(ctx context.Context, cr metav1.Object, obj client.Object) error
	DeleteObject(ctx context.Context, obj client.Object) error
}
type ApiClient struct {
	Client client.Client
	Scheme *runtime.Scheme
}

func NewApiClient(client client.Client, scheme *runtime.Scheme) *ApiClient {
	return &ApiClient{
		Client: client,
		Scheme: scheme,
	}
}

func (r *ApiClient) EnsureObject(ctx context.Context, cr metav1.Object, obj client.Object) error {
	// TODO: test how this behaves when there is f.e. CM and Deploy with same namespacedname
	log.V(2).Info("Setting controller reference", "owner", cr, "object", obj)
	if err := controllerutil.SetControllerReference(cr, obj, r.Scheme); err != nil {
		return err
	}
	log.V(2).Info("Controller reference set", "owner", cr, "object", obj)

	err := r.Client.Create(ctx, obj)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			err = r.Client.Update(ctx, obj)
			if err != nil {
				log.Error(err, "Error updating object", "Object", obj)
				return err
			}
			log.Info("Successfuly updated object", "Object", obj)
			return nil
		}
		log.Error(err, "Error creating new Object")
		return err
	}
	log.Info("Successfuly created new Object", "Object", obj)
	return nil
}

func (r *ApiClient) DeleteDeployment(ctx context.Context, cr metav1.Object, name string) error {
	d := &appsv1.Deployment{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.GetNamespace()}, d)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, d, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}

func (r *ApiClient) DeleteConfigMap(ctx context.Context, cr metav1.Object, name string) error {
	cm := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.GetNamespace()}, cm)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, cm, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}

func (r *ApiClient) DeleteService(ctx context.Context, cr metav1.Object, name string) error {
	svc := &corev1.Service{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.GetNamespace()}, svc)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, svc, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil

}

func (r *ApiClient) DeleteIngress(ctx context.Context, cr metav1.Object, name string) error {
	ingress := &networkingv1.Ingress{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: cr.GetNamespace()}, ingress)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, ingress, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}

func (r *ApiClient) DeleteObject(ctx context.Context, obj client.Object) error {
	err := r.Client.Get(ctx, types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, obj)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err == nil {
		r.Client.Delete(ctx, obj, client.PropagationPolicy(metav1.DeletePropagationForeground))
	}
	return nil
}
