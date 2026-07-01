/*
Copyright 2024.

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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/gateway"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/janekbaraniewski/kubeserial/pkg/utils/apis"
)

var devLog = logf.Log.WithName("DeviceController")

// SerialDeviceReconciler reconciles a SerialDevice object.
type SerialDeviceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	APIClient api.API
	Namespace string
	FS        utils.FileSystem
}

//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SerialDeviceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := devLog.WithName("Reconcile")
	logger.Info("Starting device reconcile", "req", req)

	defer r.updateMetrics(ctx, logger)

	device := &kubeserialv1alpha1.SerialDevice{}
	err := r.Get(ctx, req.NamespacedName, device)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "Device not found", "req", req)
			return ctrl.Result{
				Requeue: false,
			}, nil
		}
		logger.Error(err, "Failed getting device instance, will try again")
		return ctrl.Result{
			Requeue: true,
		}, nil
	}
	logger = logger.WithValues("device", device)

	err = r.EnsureConditions(ctx, device)
	if err != nil {
		logger.Error(err, "Failed ensuring conditions")
		return ctrl.Result{}, err
	}

	err = r.ValidateDeviceReady(ctx, device, req)
	if err != nil {
		logger.Error(err, "Failed device validation")
		return ctrl.Result{}, err
	}

	if device.IsAvailable() {
		return r.HandleDeviceAvailable(ctx, device, logger)
	}
	return r.HandleDeviceUnavailable(ctx, device, logger)
}

func (r *SerialDeviceReconciler) HandleDeviceAvailable(
	ctx context.Context,
	device *kubeserialv1alpha1.SerialDevice,
	_ logr.Logger,
) (reconcile.Result, error) {
	err := r.RequestGateway(ctx, device)
	if err != nil {
		//nolint:nilerr
		return ctrl.Result{}, nil
	}
	if !device.NeedsManager() {
		return ctrl.Result{}, nil
	}
	err = r.RequestManager(ctx, device)
	return reconcile.Result{}, err
}

func (r *SerialDeviceReconciler) HandleDeviceUnavailable(
	ctx context.Context,
	device *kubeserialv1alpha1.SerialDevice,
	_ logr.Logger,
) (reconcile.Result, error) {
	err := r.EnsureNoGatewayRunning(ctx, device)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !device.NeedsManager() {
		return ctrl.Result{}, nil
	}
	err = r.EnsureNoManagerRequested(ctx, device)
	return ctrl.Result{}, err
}

func (r *SerialDeviceReconciler) RequestGateway(ctx context.Context, device *kubeserialv1alpha1.SerialDevice) error {
	gatewayObjects, err := gateway.NewBuilder(device, r.FS).Build()
	if err != nil {
		return fmt.Errorf("building gateway objects: %w", err)
	}

	for _, o := range gatewayObjects {
		if err := r.APIClient.EnsureObject(ctx, device, o); err != nil {
			return err
		}
	}

	return nil
}

func (r *SerialDeviceReconciler) EnsureNoGatewayRunning(ctx context.Context, device *kubeserialv1alpha1.SerialDevice) error {
	name := apis.GatewayName(device.Name)

	if err := r.APIClient.DeleteObject(
		ctx,
		&appsv1.Deployment{ObjectMeta: v1.ObjectMeta{Name: name, Namespace: r.Namespace}},
	); err != nil {
		return err
	}

	if err := r.APIClient.DeleteObject(
		ctx,
		&corev1.ConfigMap{ObjectMeta: v1.ObjectMeta{Name: name, Namespace: r.Namespace}},
	); err != nil {
		return err
	}
	if err := r.APIClient.DeleteObject( //nolint:if-return
		ctx,
		&corev1.Service{ObjectMeta: v1.ObjectMeta{Name: name, Namespace: r.Namespace}},
	); err != nil {
		return err
	}
	return nil
}

// EnsureConditions makes sure all conditions are available in resource.
func (r *SerialDeviceReconciler) EnsureConditions(ctx context.Context, device *kubeserialv1alpha1.SerialDevice) error {
	logger := devLog.WithName("EnsureConditions")
	for _, conditionType := range []kubeserialv1alpha1.SerialDeviceConditionType{
		kubeserialv1alpha1.SerialDeviceAvailable,
		kubeserialv1alpha1.SerialDeviceReady,
		kubeserialv1alpha1.SerialDeviceFree,
	} {
		if device.GetCondition(conditionType) == nil {
			logger.Info("Condition didn't exist, creating", "conditionType", conditionType)
			device.SetCondition(kubeserialv1alpha1.SerialDeviceCondition{
				Type:   conditionType,
				Status: v1.ConditionFalse,
				Reason: "NotValidated",
			})
		}
	}
	return r.Status().Update(ctx, device)
}

// ValidateDeviceReady validates if device config is ready to be used.
func (r *SerialDeviceReconciler) ValidateDeviceReady(
	ctx context.Context,
	device *kubeserialv1alpha1.SerialDevice,
	req reconcile.Request,
) error {
	logger := devLog.WithName("ValidateDeviceReady")
	if !device.IsReady() {
		valid, err := r.ValidateDeviceManager(ctx, device, req)
		if err != nil {
			return err
		}
		if !valid {
			return nil
		}
		logger.Info("All checks passed, device ready")
		device.SetCondition(kubeserialv1alpha1.SerialDeviceCondition{
			Type:   kubeserialv1alpha1.SerialDeviceReady,
			Status: v1.ConditionTrue,
			Reason: "AllChecksPassed",
		})
		if err := r.Status().Update(ctx, device); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDeviceManager validates if device manaager config is valid and upadates device state in case it's not.
func (r *SerialDeviceReconciler) ValidateDeviceManager(
	ctx context.Context,
	device *kubeserialv1alpha1.SerialDevice,
	req reconcile.Request,
) (bool, error) {
	if !device.NeedsManager() {
		return true, nil
	}
	if !r.ManagerIsAvailable(ctx, device, req) {
		device.SetCondition(kubeserialv1alpha1.SerialDeviceCondition{
			Type:   kubeserialv1alpha1.SerialDeviceReady,
			Status: v1.ConditionFalse,
			Reason: "ManagerNotAvailable",
		})
		if err := r.Status().Update(ctx, device); err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// ManagerIsAvailable checks if Manager object referenced by SerialDevice is available in the cluster.
func (r *SerialDeviceReconciler) ManagerIsAvailable(
	ctx context.Context,
	device *kubeserialv1alpha1.SerialDevice,
	req ctrl.Request,
) bool {
	logger := devLog.WithName("ManagerIsAvailable")
	manager := &kubeserialv1alpha1.Manager{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      device.Spec.Manager,
		Namespace: req.Namespace,
	}, manager)
	if err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		logger.Error(err, "Unknown error")
	}

	return true
}

// RequestManager create ManagerScheduleRequest for device.
func (r *SerialDeviceReconciler) RequestManager(ctx context.Context, device *kubeserialv1alpha1.SerialDevice) error {
	logger := devLog.WithName("RequestManager")
	if !device.NeedsManager() {
		logger.V(3).Info("Device doesn't need manager, returning")
		return nil
	}
	request := &kubeserialv1alpha1.ManagerScheduleRequest{
		ObjectMeta: v1.ObjectMeta{
			Name:      apis.ScheduleRequestName(device.Name, device.Spec.Manager),
			Namespace: device.Namespace,
		},
		Spec: kubeserialv1alpha1.ManagerScheduleRequestSpec{
			Device:  device.Name,
			Manager: device.Spec.Manager,
		},
	}
	if err := controllerutil.SetControllerReference(device, request, r.Scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}
	err := r.Create(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

// EnsureNoManagerRequested makes sure there is no ManagerScheduleRequest for device.
func (r *SerialDeviceReconciler) EnsureNoManagerRequested(ctx context.Context, device *kubeserialv1alpha1.SerialDevice) error {
	if !device.NeedsManager() {
		return nil
	}
	request := &kubeserialv1alpha1.ManagerScheduleRequest{}
	if err := r.Get(ctx, types.NamespacedName{
		Name:      apis.ScheduleRequestName(device.Name, device.Spec.Manager),
		Namespace: device.Namespace,
	}, request); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return r.Delete(ctx, request)
}

// updateMetrics refreshes the device gauges from the full set of SerialDevices.
func (r *SerialDeviceReconciler) updateMetrics(ctx context.Context, logger logr.Logger) {
	list := &kubeserialv1alpha1.SerialDeviceList{}
	if err := r.List(ctx, list); err != nil {
		logger.Error(err, "Failed listing devices for metrics")
		return
	}

	var available, ready float64
	for i := range list.Items {
		if list.Items[i].IsAvailable() {
			available++
		}
		if list.Items[i].IsReady() {
			ready++
		}
	}

	Devices.Set(float64(len(list.Items)))
	AvailableDevices.Set(available)
	ReadyDevices.Set(ready)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SerialDeviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.SerialDevice{}).
		Complete(r)
}
