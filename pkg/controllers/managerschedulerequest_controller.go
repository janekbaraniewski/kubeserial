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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"github.com/janekbaraniewski/kubeserial/pkg/managers"
)

var msrcLog = logf.Log.WithName("ManagerScheduleRequestController")

// ManagerScheduleRequestReconciler reconciles a ManagerScheduleRequest object
type ManagerScheduleRequestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubeserial.app.kubeserial.com,resources=managerschedulerequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ManagerScheduleRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := msrcLog.WithName("Reconcile")
	instance := &kubeserialv1alpha1.ManagerScheduleRequest{}
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
	logger = logger.WithValues("request", instance)

	manager := &kubeserialv1alpha1.Manager{}
	err = r.Client.Get(ctx, types.NamespacedName{
		Name:      instance.Spec.Manager,
		Namespace: req.Namespace,
	}, manager)
	if err != nil {
		// TODO: handle missing manager spec
		return ctrl.Result{}, nil
	}
	logger = logger.WithValues("manager", manager)
	logger.Info("Got manager, starting ReconcileManager")
	r.ReconcileManager(ctx, instance, manager, req)

	return ctrl.Result{}, nil
}

// ReconcileManager
func (r *ManagerScheduleRequestReconciler) ReconcileManager(ctx context.Context, instance *kubeserialv1alpha1.ManagerScheduleRequest, mgr *kubeserialv1alpha1.Manager, req ctrl.Request) error {
	apiClient := api.ApiClient{
		Client: r.Client,
		Scheme: r.Scheme,
	}
	if err := managers.ScheduleFromCRD(ctx, instance, mgr, &apiClient); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManagerScheduleRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubeserialv1alpha1.ManagerScheduleRequest{}).
		Complete(r)
}
