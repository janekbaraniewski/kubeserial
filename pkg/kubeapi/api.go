package kubeapi

import (
	"context"

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
	log.V(2).Info("Controller reference set", "ownerName", cr.GetName(), "objectName", obj.GetName())

	err := r.Client.Create(ctx, obj)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			err = r.Client.Update(ctx, obj)
			if err != nil {
				log.Error(err, "Error updating object", "ObjectName", obj.GetName())
				return err
			}
			log.Info("Successfuly updated object", "ObjectName", obj.GetName())
			return nil
		}
		log.Error(err, "Error creating new Object")
		return err
	}
	log.Info("Successfuly created new Object", "ObjectName", obj.GetName())
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
