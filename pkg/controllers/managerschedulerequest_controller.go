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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

var msrcLog = logf.Log.WithName("ManagerScheduleRequestController")

// ManagerScheduleRequestReconciler reconciles a ManagerScheduleRequest object.
type ManagerScheduleRequestReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Namespace string
	FS        utils.FileSystem
	APIClient api.API
}

//nolint
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ManagerScheduleRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := msrcLog.WithName("Reconcile")
	instance := &kubeserialv1alpha1.ManagerScheduleRequest{}
	err := r.Get(ctx, req.NamespacedName, instance)
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
	logger = logger.WithValues("request", instance)

	manager := &kubeserialv1alpha1.Manager{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      instance.Spec.Manager,
		Namespace: req.Namespace,
	}, manager)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Referenced manager not found, marking request unfulfilled and requeueing",
				"manager", instance.Spec.Manager)
			if statusErr := r.setFulfilled(ctx, instance, false); statusErr != nil {
				return ctrl.Result{}, statusErr
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, fmt.Errorf("getting manager %q: %w", instance.Spec.Manager, err)
	}
	logger = logger.WithValues("manager", manager)
	logger.Info("Got manager, starting ReconcileManager")
	if err := r.ReconcileManager(ctx, instance, manager, req); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.setFulfilled(ctx, instance, true); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// setFulfilled updates the request's Fulfilled status if it changed.
func (r *ManagerScheduleRequestReconciler) setFulfilled(
	ctx context.Context,
	instance *kubeserialv1alpha1.ManagerScheduleRequest,
	fulfilled bool,
) error {
	if instance.Status.Fulfilled == fulfilled {
		return nil
	}
	instance.Status.Fulfilled = fulfilled
	if err := r.Status().Update(ctx, instance); err != nil {
		return fmt.Errorf("updating ManagerScheduleRequest status: %w", err)
	}
	return nil
}

// ReconcileManager.
func (r *ManagerScheduleRequestReconciler) ReconcileManager(
	ctx context.Context,
	instance *kubeserialv1alpha1.ManagerScheduleRequest,
	mgr *kubeserialv1alpha1.Manager,
	_ ctrl.Request,
) error {
	return managers.Schedule(ctx, r.FS, instance, mgr, r.Namespace, r.APIClient)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManagerScheduleRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.ManagerScheduleRequest{}).
		Complete(r)
}
