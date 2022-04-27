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
	devLog.Info("Starting device reconcile", "req", req)

	instance := &kubeserialv1alpha1.Device{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			devLog.Error(err, "Device not found", "req", req)
			return ctrl.Result{
				Requeue: false,
			}, nil
		}
		devLog.Error(err, "Failed getting device instance, will try again")
		return ctrl.Result{
			Requeue: true,
		}, nil
	}
	log := devLog.WithValues("device", instance)

	err = r.EnsureConditions(ctx, instance)
	if err != nil {
		log.Error(err, "Failed ensuring conditions")
		return ctrl.Result{}, err
	}

	err = r.ValidateDeviceReady(ctx, instance, req)
	if err != nil {
		log.Error(err, "Failed device validation")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// EnsureConditions makes sure all conditions are available in resource
func (r *DeviceReconciler) EnsureConditions(ctx context.Context, instance *kubeserialv1alpha1.Device) error {
	for _, conditionType := range []kubeserialv1alpha1.DeviceConditionType{
		kubeserialv1alpha1.DeviceAvailable,
		kubeserialv1alpha1.DeviceReady,
	} {
		if utils.GetCondition(instance.Status.Conditions, conditionType) == nil {
			log.Info("Condition didn't exist, creating", "conditionType", conditionType)
			utils.SetDeviceCondition(&instance.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
				Type:   conditionType,
				Status: v1.ConditionUnknown,
				Reason: "NotValidated",
			})
		}
	}
	if err := r.Client.Update(ctx, instance); err != nil {
		return err
	}
	return nil
}

// ValidateDeviceReady validates if device config is ready to be used
func (r *DeviceReconciler) ValidateDeviceReady(ctx context.Context, instance *kubeserialv1alpha1.Device, req reconcile.Request) error {
	readyCondition := utils.GetCondition(instance.Status.Conditions, v1alpha1.DeviceReady)
	if readyCondition.Status != v1.ConditionTrue {
		err := r.ValidateDeviceManager(ctx, instance, req)
		if err != nil {
			return err
		}
		log.Info("All checks passed, device ready")
		utils.SetDeviceCondition(&instance.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
			Type:   kubeserialv1alpha1.DeviceReady,
			Status: v1.ConditionTrue,
			Reason: "AllChecksPassed",
		})
		if err := r.Client.Update(ctx, instance); err != nil {
			return err
		}
	}
	return nil
}

// ValidateDeviceManager validates if device manaager config is valid and upadates device state in case it's not
func (r *DeviceReconciler) ValidateDeviceManager(ctx context.Context, instance *kubeserialv1alpha1.Device, req reconcile.Request) error {
	if !r.ManagerIsAvailable(ctx, instance, req) {
		utils.SetDeviceCondition(&instance.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
			Type:   kubeserialv1alpha1.DeviceReady,
			Status: v1.ConditionFalse,
			Reason: "ManagerNotAvailable",
		})
		if err := r.Client.Update(ctx, instance); err != nil {
			return err
		}
	}
	return nil
}

// ManagerIsAvailable checks if Manager object referenced by Device is available in the cluster
func (r *DeviceReconciler) ManagerIsAvailable(ctx context.Context, device *kubeserialv1alpha1.Device, req ctrl.Request) bool {
	log := devLog.WithName("ManagerIsAvailable")
	manager := &kubeserialv1alpha1.Manager{}

	err := r.Client.Get(ctx, types.NamespacedName{
		Name:      device.Spec.Manager,
		Namespace: req.Namespace,
	}, manager)

	if err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		log.Error(err, "Unknown error")
	}

	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeviceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.Device{}).
		Complete(r)
}
