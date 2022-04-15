package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
)

func GetTestSetup() (KubeSerialReconciler, client.WithWatch, api.FakeApiClient) {
	fakeClient := runtimefake.NewClientBuilder().Build()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))

	reconciler := KubeSerialReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	apiClient := api.NewFakeApiClient()

	return reconciler, fakeClient, apiClient
}

func TestReconcileGateway_ScheduleForNewDevice(t *testing.T) {
	reconciler, fakeClient, apiClient := GetTestSetup()
	cr := &kubeserialv1alpha1.KubeSerial{
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

	deviceConf := CreateDeviceConfig(cr)[0]
	deviceConf.Data["available"] = "true"
	fakeClient.Create(context.TODO(), deviceConf)

	reconciler.ReconcileGateway(context.TODO(), cr, &apiClient)
	assert.Equal(t,
		[]string{
			"EnsureConfigMap",
			"EnsureDeployment",
			"EnsureService",
		}, apiClient.Operations)
}

func TestReconcileGateway_DeleteForMissingDevice(t *testing.T) {
	reconciler, fakeClient, apiClient := GetTestSetup()
	cr := &kubeserialv1alpha1.KubeSerial{
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

	deviceConf := CreateDeviceConfig(cr)[0]
	deviceConf.Data["available"] = "false"
	fakeClient.Create(context.TODO(), deviceConf)

	reconciler.ReconcileGateway(context.TODO(), cr, &apiClient)
	assert.Equal(t,
		[]string{
			"DeleteDeployment",
			"DeleteConfigMap",
			"DeleteService",
		}, apiClient.Operations)
}
