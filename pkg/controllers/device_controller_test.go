package controllers

import (
	"context"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	controllerruntime "sigs.k8s.io/controller-runtime"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestDeviceReconciler_Reconcile(t *testing.T) {
	deviceName := types.NamespacedName{
		Name:      "test-device",
		Namespace: "test-ns",
	}

	device := &v1alpha1.Device{
		ObjectMeta: v1.ObjectMeta{
			Name:      deviceName.Name,
			Namespace: deviceName.Namespace,
		},
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

	err := fakeClient.Create(context.TODO(), device)

	assert.Equal(t, nil, err)

	deviceReconciler := DeviceReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

	assert.Equal(t, nil, err)
	assert.Equal(t, false, result.Requeue)
}
