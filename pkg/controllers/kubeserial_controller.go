/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	apiclient "github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
)

var log = logf.Log.WithName("KubeSerialController")

// KubeSerialReconciler reconciles a KubeSerial object
type KubeSerialReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KubeSerial object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *KubeSerialReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithName("Reconcile")

	reqLogger.Info("Reconciling KubeSerial")

	instance := &kubeserialv1alpha1.KubeSerial{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	apiClient := apiclient.ApiClient{
		Client: r.Client,
		Scheme: r.Scheme,
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

func (r *KubeSerialReconciler) reconcileDevicesConfig(cr *appv1alpha1.KubeSerial, api *apiclient.ApiClient) error {
	logger := log.WithName("reconcileDevicesConfig")
	deviceConfs := createDeviceConfig(cr)

	for _, deviceConf := range deviceConfs {
		if err := controllerutil.SetControllerReference(cr, deviceConf, r.Scheme); err != nil {
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

func (r *KubeSerialReconciler) GetDeviceState(p *appv1alpha1.Device, cr *appv1alpha1.KubeSerial) (*corev1.ConfigMap, error) {
	logger := log.WithName("GetDevicesState")

	found := &corev1.ConfigMap{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: strings.ToLower(cr.Name + "-" + p.Name), Namespace: cr.Namespace}, found)
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
		labels := map[string]string{
			"app":    cr.Name,
			"device": device.Name,
			"type":   "DeviceState",
		}

		confs = append(confs, &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(cr.Name + "-" + device.Name),
				Namespace: cr.Namespace,
				Labels:    labels,
			},
			Data: map[string]string{
				"available": "false",
				"node":      "",
			},
		})
	}
	return confs
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeSerialReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.KubeSerial{}).
		Complete(r)
}