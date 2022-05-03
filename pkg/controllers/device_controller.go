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

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

var devLog = logf.Log.WithName("DeviceController")

// DeviceReconciler reconciles a Device object
type DeviceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=devices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DeviceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := devLog.WithName("Reconcile")
	logger.Info("Starting device reconcile", "req", req)

	device := &kubeserialv1alpha1.Device{}
	err := r.Client.Get(ctx, req.NamespacedName, device)
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

	if !device.NeedsManager() {
		return ctrl.Result{}, nil
	}

	availableCondition := utils.GetCondition(device.Status.Conditions, v1alpha1.DeviceAvailable)
	if availableCondition.Status == v1.ConditionTrue {
		r.RequestManager(ctx, device)
	} else {
		r.EnsureNoManagerRequested(ctx, device)
	}

	return ctrl.Result{}, nil
}

// EnsureConditions makes sure all conditions are available in resource
func (r *DeviceReconciler) EnsureConditions(ctx context.Context, device *kubeserialv1alpha1.Device) error {
	logger := devLog.WithName("EnsureConditions")
	for _, condition := range []struct {
		Type     kubeserialv1alpha1.DeviceConditionType
		IsNeeded func() bool
	}{
		{
			Type:     kubeserialv1alpha1.DeviceAvailable,
			IsNeeded: func() bool { return true },
		},
		{
			Type:     kubeserialv1alpha1.DeviceReady,
			IsNeeded: device.NeedsManager,
		},
	} {
		if !condition.IsNeeded() {
			continue
		}
		if utils.GetCondition(device.Status.Conditions, condition.Type) == nil {
			logger.Info("Condition didn't exist, creating", "conditionType", condition.Type)
			utils.SetDeviceCondition(&device.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
				Type:   condition.Type,
				Status: v1.ConditionUnknown,
				Reason: "NotValidated",
			})
		}
	}
	if err := r.Client.Status().Update(ctx, device); err != nil {
		return err
	}
	return nil
}

// ValidateDeviceReady validates if device config is ready to be used
func (r *DeviceReconciler) ValidateDeviceReady(ctx context.Context, device *kubeserialv1alpha1.Device, req reconcile.Request) error {
	logger := devLog.WithName("ValidateDeviceReady")
	readyCondition := utils.GetCondition(device.Status.Conditions, v1alpha1.DeviceReady)
	if readyCondition.Status != v1.ConditionTrue {
		valid, err := r.ValidateDeviceManager(ctx, device, req)
		if err != nil {
			return err
		}
		if !valid {
			return nil
		}
		logger.Info("All checks passed, device ready")
		utils.SetDeviceCondition(&device.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
			Type:   kubeserialv1alpha1.DeviceReady,
			Status: v1.ConditionTrue,
			Reason: "AllChecksPassed",
		})
		if err := r.Client.Status().Update(ctx, device); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDeviceManager validates if device manaager config is valid and upadates device state in case it's not
func (r *DeviceReconciler) ValidateDeviceManager(ctx context.Context, device *kubeserialv1alpha1.Device, req reconcile.Request) (bool, error) {
	if !device.NeedsManager() {
		return true, nil
	}
	if !r.ManagerIsAvailable(ctx, device, req) {
		utils.SetDeviceCondition(&device.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
			Type:   kubeserialv1alpha1.DeviceReady,
			Status: v1.ConditionFalse,
			Reason: "ManagerNotAvailable",
		})
		if err := r.Client.Status().Update(ctx, device); err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// ManagerIsAvailable checks if Manager object referenced by Device is available in the cluster
func (r *DeviceReconciler) ManagerIsAvailable(ctx context.Context, device *kubeserialv1alpha1.Device, req ctrl.Request) bool {
	logger := devLog.WithName("ManagerIsAvailable")
	manager := &kubeserialv1alpha1.Manager{}

	err := r.Client.Get(ctx, types.NamespacedName{
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

// RequestManager create ManagerScheduleRequest for device
func (r *DeviceReconciler) RequestManager(ctx context.Context, device *kubeserialv1alpha1.Device) error {
	request := &kubeserialv1alpha1.ManagerScheduleRequest{
		ObjectMeta: v1.ObjectMeta{
			Name:      device.Name + "-" + device.Spec.Manager,
			Namespace: device.Namespace,
		},
		Spec: kubeserialv1alpha1.ManagerScheduleRequestSpec{
			Device:  device.Name,
			Manager: device.Spec.Manager,
		},
	}
	err := r.Client.Create(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

// EnsureNoManagerRequested makes sure there is no ManagerScheduleRequest for device
func (r *DeviceReconciler) EnsureNoManagerRequested(ctx context.Context, device *kubeserialv1alpha1.Device) error {
	request := &kubeserialv1alpha1.ManagerScheduleRequest{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Name:      device.Name + "-" + device.Spec.Manager,
		Namespace: device.Namespace,
	}, request); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.Client.Delete(ctx, request); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.Device{}).
		Complete(r)
}
