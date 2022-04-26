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
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

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

	log.Info("Ensuring conditions exist")
	// Ensure all conditions are available
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
	log.Info("Updating device status")
	if err := r.Client.Update(ctx, instance); err != nil {
		log.Error(err, "Failed updating device status")
		return ctrl.Result{}, err
	}

	if !r.ManagerIsAvailable(ctx, instance, req) {
		log.Info("Manager for device is unavailable. Will try again in 1 minute.")
		utils.SetDeviceCondition(&instance.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
			Type:   kubeserialv1alpha1.DeviceReady,
			Status: v1.ConditionFalse,
			Reason: "ManagerNotAvailable",
		})
		if err := r.Client.Update(ctx, instance); err != nil {
			log.Error(err, "Failed updating device status")
		}
		return ctrl.Result{
			RequeueAfter: 1 * time.Minute,
		}, nil
	}

	log.Info("All checks passed, device ready")
	utils.SetDeviceCondition(&instance.Status.Conditions, kubeserialv1alpha1.DeviceCondition{
		Type:   kubeserialv1alpha1.DeviceReady,
		Status: v1.ConditionTrue,
		Reason: "AllChecksPassed",
	})
	if err := r.Client.Update(ctx, instance); err != nil {
		log.Error(err, "Failed updating device status")
	}
	return ctrl.Result{}, nil
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
