package kubeserial

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"

	corev1 "k8s.io/api/core/v1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

)

func (r *ReconcileKubeSerial) ReconcileMonitor(cr *appv1alpha1.KubeSerial) error {
	if err := r.reconcileConfigMap(cr); err != nil {
		return err
	}

	if err := r.reconcileDaemonSet(cr); err != nil {
		return err
	}

	return nil
}

func (r *ReconcileKubeSerial) reconcileConfigMap(cr *appv1alpha1.KubeSerial) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)
	conf 	:= monitor.CreateConfigMap(cr)
	if err := controllerutil.SetControllerReference(cr, conf, r.scheme); err != nil {
			return err
	}

	found := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: conf.Name, Namespace: conf.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.client.Create(context.TODO(), conf)
		if err != nil {
			logger.Info(err.Error())
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileKubeSerial) reconcileDaemonSet(cr *appv1alpha1.KubeSerial) error {
	logger := log.WithValues("KubeSerial.Namespace", cr.Namespace, "KubeSerial.Name", cr.Name)
	monitorDaemon 	:= monitor.CreateDaemonSet(cr)

	if err := controllerutil.SetControllerReference(cr, monitorDaemon, r.scheme); err != nil {
		logger.Info("Can't set reference")
		return err
	}

	found := &v1beta2.DaemonSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: monitorDaemon.Name, Namespace: monitorDaemon.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		err = r.client.Create(context.TODO(), monitorDaemon)
		if err != nil {
			logger.Info(err.Error())
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
