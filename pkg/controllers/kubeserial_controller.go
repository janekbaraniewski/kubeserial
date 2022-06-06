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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

var ksLog = logf.Log.WithName("KubeSerialController")

// KubeSerialReconciler reconciles a KubeSerial object
type KubeSerialReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	DeviceMonitorVersion string
	FS                   utils.FileSystem
	APIClient            api.API
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

	if err := r.ReconcileMonitor(ctx, instance, r.DeviceMonitorVersion); err != nil {
		reqLogger.Info("ReconcileMonitor fail")
		return reconcile.Result{}, err
	}

	return reconcile.Result{
		RequeueAfter: time.Second * 5,
	}, nil
}

func (r *KubeSerialReconciler) ReconcileMonitor(ctx context.Context, cr *appv1alpha1.KubeSerial, monitorVersion string) error {
	cm, err := monitor.CreateConfigMap(r.FS, cr.Spec.SerialDevices)
	if err != nil {
		return err
	}

	monitorDaemon, err := monitor.CreateDaemonSet(r.FS)

	if err != nil {
		return err
	}

	if err := r.APIClient.EnsureObject(ctx, cr, cm); err != nil {
		return err
	}

	if err := r.APIClient.EnsureObject(ctx, cr, monitorDaemon); err != nil {
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubeSerialReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.KubeSerial{}).
		Complete(r)
}
