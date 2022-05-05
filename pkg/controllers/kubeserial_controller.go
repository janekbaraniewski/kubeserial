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
	"time"

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
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
)

var ksLog = logf.Log.WithName("KubeSerialController")

// KubeSerialReconciler reconciles a KubeSerial object
type KubeSerialReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	DeviceMonitorVersion string
}

//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.kubeserial.com,resources=kubeserials/finalizers,verbs=update
func (r *KubeSerialReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := ksLog.WithName("Reconcile")

	reqLogger.Info("Reconciling KubeSerial")

	instance := &appv1alpha1.KubeSerial{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	apiClient := api.ApiClient{
		Client: r.Client,
		Scheme: r.Scheme,
	}

	if err := r.reconcileDevicesConfig(ctx, instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileConfig fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileMonitor(ctx, instance, &apiClient, r.DeviceMonitorVersion); err != nil {
		reqLogger.Info("ReconcileMonitor fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileGateway(ctx, instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileGateway fail")
		return reconcile.Result{}, err
	}

	if err := r.ReconcileManagers(ctx, instance, &apiClient); err != nil {
		reqLogger.Info("ReconcileManager fail")
		return reconcile.Result{}, err
	}

	return reconcile.Result{
		RequeueAfter: time.Second * 5,
	}, nil
}

func (r *KubeSerialReconciler) ReconcileManagers(ctx context.Context, cr *appv1alpha1.KubeSerial, api api.API) error {
	for _, device := range cr.Spec.Devices {
		stateCM, err := r.GetDeviceState(ctx, &device, cr)
		if err != nil {
			return err
		}
		manager := managers.Available[device.Manager]
		if stateCM.Data["available"] == "true" {
			if err := manager.Schedule(ctx, cr, &device, api); err != nil {
				return err
			}
		} else {
			if err := manager.Delete(ctx, cr, &device, api); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *KubeSerialReconciler) ReconcileMonitor(ctx context.Context, cr *appv1alpha1.KubeSerial, api api.API, monitorVersion string) error {
	conf := monitor.CreateConfigMap(cr)
	monitorDaemon := monitor.CreateDaemonSet(cr, monitorVersion)

	if err := api.EnsureConfigMap(ctx, cr, conf); err != nil {
		return err
	}

	if err := api.EnsureDaemonSet(ctx, cr, monitorDaemon); err != nil {
		return err
	}

	return nil
}

func (r *KubeSerialReconciler) reconcileDevicesConfig(ctx context.Context, cr *appv1alpha1.KubeSerial, api api.API) error {
	logger := ksLog.WithName("reconcileDevicesConfig")
	deviceConfs := CreateDeviceConfig(cr)

	for _, deviceConf := range deviceConfs {
		if err := controllerutil.SetControllerReference(cr, deviceConf, r.Scheme); err != nil {
			logger.Info("Can't set reference")
			return err
		}
	}

	for _, deviceConf := range deviceConfs {
		if err := api.EnsureConfigMap(ctx, cr, deviceConf); err != nil {
			return err
		}
	}

	return nil
}

func (r *KubeSerialReconciler) GetDeviceState(ctx context.Context, p *appv1alpha1.Device_2, cr *appv1alpha1.KubeSerial) (*corev1.ConfigMap, error) {
	logger := ksLog.WithName("GetDevicesState").WithValues("Device", p.Name)
	logger.Info("Fetching device state")
	found := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: strings.ToLower(cr.Name + "-" + p.Name), Namespace: cr.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Can't get device state")
		return nil, err
	} else if err != nil {
		return nil, err
	}
	logger.Info("Got state CM", "ConfigMap", found.Data)
	return found, nil
}

func CreateDeviceConfig(cr *appv1alpha1.KubeSerial) []*corev1.ConfigMap { // TODO: move to separate module
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
		For(&appv1alpha1.KubeSerial{}).
		Complete(r)
}
