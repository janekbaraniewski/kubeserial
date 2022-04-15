package controllers

import (
	"context"
	"testing"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func getCR() *kubeserialv1alpha1.KubeSerial {
	return &kubeserialv1alpha1.KubeSerial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubeserialtest",
			Namespace: "test-namespace",
		},
		Spec: kubeserialv1alpha1.KubeSerialSpec{
			Devices: []kubeserialv1alpha1.Device{
				{
					Name:      "testDevice",
					IdVendor:  "0",
					IdProduct: "1",
					Manager:   "testManager",
					Subsystem: "tty",
				},
			},
		},
	}
}

func TestCreateDeviceConfig(t *testing.T) {
	cr := getCR()

	devices := CreateDeviceConfig(cr)

	desired := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubeserialtest-testdevice",
			Namespace: "test-namespace",
			Labels: map[string]string{
				"app":    "kubeserialtest",
				"device": "testDevice",
				"type":   "DeviceState",
			},
		},
		Data: map[string]string{
			"available": "false",
			"node":      "",
		},
	}

	assert.Equal(t, 1, len(devices))
	assert.Equal(t, &desired, devices[0])
}

func TestGetDeviceState(t *testing.T) {
	fakeClient := runtimefake.NewClientBuilder().Build()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))

	reconciler := KubeSerialReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	testCases := []struct {
		Name          string
		ObjToCreate   client.Object
		ExpectedState *corev1.ConfigMap
		ValidateError func(error) bool
		Device        kubeserialv1alpha1.Device
		KS            *kubeserialv1alpha1.KubeSerial
	}{
		{
			Name:          "cant_find_device",
			ObjToCreate:   nil,
			ExpectedState: nil,
			ValidateError: func(err error) bool { return errors.IsNotFound(err) },
			Device:        kubeserialv1alpha1.Device{},
			KS:            getCR(),
		},
		{
			Name: "success",
			ObjToCreate: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubeserialtest-testdevice",
					Namespace: "test-namespace",
				},
			},
			ExpectedState: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "kubeserialtest-testdevice",
					Namespace:       "test-namespace",
					ResourceVersion: "1",
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
			},
			ValidateError: func(err error) bool { return err == nil },
			Device:        kubeserialv1alpha1.Device{Name: "testdevice"},
			KS:            getCR(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			if testCase.ObjToCreate != nil {
				fakeClient.Create(context.TODO(), testCase.ObjToCreate)
			}

			state, err := reconciler.GetDeviceState(context.TODO(), &testCase.Device, testCase.KS)
			assert.Equal(t, testCase.ExpectedState, state)
			assert.Equal(t, true, testCase.ValidateError(err))
		})
	}
}

func TestReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

	reconciler := KubeSerialReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	reconcileReq := reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "test-namespace", Name: "kubeserialtest"}}

	result, err := reconciler.Reconcile(context.TODO(), reconcileReq)
	assert.Equal(t, nil, err)
	assert.Equal(t, reconcile.Result{}, result)
}
