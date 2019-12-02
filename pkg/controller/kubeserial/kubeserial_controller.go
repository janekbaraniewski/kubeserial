package kubeserial

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	apiclient "github.com/janekbaraniewski/kubeserial/pkg/controller/api"
)

var log = logf.Log.WithName("controller_kubeserial")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KubeSerial Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKubeSerial{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kubeserial-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource KubeSerial
	err = c.Watch(&source.Kind{Type: &appv1alpha1.KubeSerial{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1beta2.DaemonSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1beta2.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.KubeSerial{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKubeSerial implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKubeSerial{}

// ReconcileKubeSerial reconciles a KubeSerial object
type ReconcileKubeSerial struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a KubeSerial object and makes changes based on the state read
// and what is in the KubeSerial.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKubeSerial) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KubeSerial")

	instance := &appv1alpha1.KubeSerial{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	apiClient := apiclient.ApiClient{
		Client:		r.client,
		Scheme:		r.scheme,
	}

	if err := r.reconcileDevicesConfig(instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileConfig fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileMonitor(instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileMonitor fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileGateway(instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileGateway fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileManagers(instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileManager fail")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileKubeSerial) reconcileDevicesConfig(cr *appv1alpha1.KubeSerial, api *apiclient.ApiClient) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)
	deviceConfs 	:= createDeviceConfig(cr)

	for _, deviceConf := range deviceConfs{
		if err := controllerutil.SetControllerReference(cr, deviceConf, r.scheme); err != nil {
				logger.Info("Can't set reference")
				return err
		}
	}

	for _, deviceConf := range deviceConfs {
		if err := api.EnsureConfigMap(cr, deviceConf); err != nil {
			return err
		}
	}

	return nil
}

func (r *ReconcileKubeSerial) GetDeviceState(p *appv1alpha1.Device, cr *appv1alpha1.KubeSerial) (*corev1.ConfigMap, error) {
	logger := log.WithValues("Namespace", cr.Namespace, "Name", cr.Name)

	found := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: cr.Name + "-" + p.Name, Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Can't get device state")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return found, nil
}

func createDeviceConfig(cr *appv1alpha1.KubeSerial) []*corev1.ConfigMap { // TODO: move to separate module
	confs := []*corev1.ConfigMap{}
	for _, device := range cr.Spec.Devices {
		labels := map[string]string {
			"app":		cr.Name,
			"device":	device.Name,
			"type":		"DeviceState",
		}

		confs = append(confs, &corev1.ConfigMap {
			ObjectMeta:	metav1.ObjectMeta {
				Name:		cr.Name + "-" + device.Name,
				Namespace:	cr.Namespace,
				Labels:		labels,
			},
			Data:		map[string]string {
				"available": "false",
				"node":		 "",
			},
		})
	}
	return confs
}
