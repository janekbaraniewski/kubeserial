package controllers

import (
	"context"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func newTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.Install(scheme))
	return scheme
}

func newTestDevice() *v1alpha1.SerialDevice {
	return &v1alpha1.SerialDevice{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-device",
			Namespace: "test-ns",
		},
		Spec: v1alpha1.SerialDeviceSpec{
			Manager:   "test-manager",
			IDVendor:  "123",
			IDProduct: "456",
			Name:      "test-device",
		},
	}
}

func newTestManager() *v1alpha1.Manager {
	return &v1alpha1.Manager{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-manager",
			Namespace: "test-ns",
		},
		Spec: v1alpha1.ManagerSpec{
			Image: v1alpha1.Image{
				Repository: "test-image",
				Tag:        "latest",
			},
			Config:     "test-config",
			ConfigPath: "/home/config.yaml",
			RunCmd:     "./test-manager",
		},
	}
}

func TestDeviceReconciler_Reconcile(t *testing.T) {
	t.Parallel()
	deviceName := types.NamespacedName{
		Name:      "test-device",
		Namespace: "test-ns",
	}

	newReconciler := func(scheme *runtime.Scheme, objs ...client.Object) (SerialDeviceReconciler, client.Client) {
		fs := utils.NewInMemoryFS()
		AddGatewaySpecFilesToFilesystem(t, fs)
		fakeClient := runtimefake.NewClientBuilder().
			WithScheme(scheme).
			WithStatusSubresource(&v1alpha1.SerialDevice{}).
			WithObjects(objs...).
			Build()
		return SerialDeviceReconciler{
			Client:    fakeClient,
			Scheme:    scheme,
			APIClient: kubeapi.NewFakeAPIClient(),
			FS:        fs,
		}, fakeClient
	}

	t.Run("device-new-manager-not-available", func(t *testing.T) {
		t.Parallel()
		scheme := newTestScheme(t)
		reconciler, fakeClient := newReconciler(scheme, newTestDevice())

		result, err := reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})
		require.NoError(t, err)
		assert.Equal(t, controllerruntime.Result{}, result)

		foundDevice := &v1alpha1.SerialDevice{}
		require.NoError(t, fakeClient.Get(context.TODO(), deviceName, foundDevice))
		assert.Len(t, foundDevice.Status.Conditions, 3)

		availableCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceAvailable)
		require.NotNil(t, availableCondition)
		assert.Equal(t, v1.ConditionFalse, availableCondition.Status)
		assert.Equal(t, "NotValidated", availableCondition.Reason)

		readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
		require.NotNil(t, readyCondition)
		assert.Equal(t, v1.ConditionFalse, readyCondition.Status)
		assert.Equal(t, "ManagerNotAvailable", readyCondition.Reason)
	})

	t.Run("device-new-manager-available", func(t *testing.T) {
		t.Parallel()
		scheme := newTestScheme(t)
		reconciler, fakeClient := newReconciler(scheme, newTestDevice(), newTestManager())

		result, err := reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})
		require.NoError(t, err)
		assert.Equal(t, controllerruntime.Result{}, result)

		foundDevice := &v1alpha1.SerialDevice{}
		require.NoError(t, fakeClient.Get(context.TODO(), deviceName, foundDevice))

		availableCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceAvailable)
		require.NotNil(t, availableCondition)
		assert.Equal(t, v1.ConditionFalse, availableCondition.Status)
		assert.Equal(t, "NotValidated", availableCondition.Reason)

		readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
		require.NotNil(t, readyCondition)
		assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
		assert.Equal(t, "AllChecksPassed", readyCondition.Reason)
	})

	t.Run("device-ready", func(t *testing.T) {
		t.Parallel()
		scheme := newTestScheme(t)
		device := newTestDevice()
		device.Status.Conditions = append(device.Status.Conditions, v1alpha1.SerialDeviceCondition{
			Type:   v1alpha1.SerialDeviceAvailable,
			Status: v1.ConditionTrue,
		})
		reconciler, fakeClient := newReconciler(scheme, device, newTestManager())

		result, err := reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})
		require.NoError(t, err)
		assert.Equal(t, controllerruntime.Result{}, result)

		foundDevice := &v1alpha1.SerialDevice{}
		require.NoError(t, fakeClient.Get(context.TODO(), deviceName, foundDevice))

		readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
		require.NotNil(t, readyCondition)
		assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
		assert.Equal(t, "AllChecksPassed", readyCondition.Reason)

		// Device controller only creates the request; the MSR controller fulfills it.
		foundRequest := &v1alpha1.ManagerScheduleRequest{}
		require.NoError(t, fakeClient.Get(context.TODO(), types.NamespacedName{
			Name:      device.Name + "-" + device.Spec.Manager,
			Namespace: device.Namespace,
		}, foundRequest))
		assert.False(t, foundRequest.Status.Fulfilled)
	})

	t.Run("device-not-found", func(t *testing.T) {
		t.Parallel()
		scheme := newTestScheme(t)
		reconciler, _ := newReconciler(scheme)

		result, err := reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})
		require.NoError(t, err)
		assert.Equal(t, controllerruntime.Result{}, result)
	})
}

func AddGatewaySpecFilesToFilesystem(t *testing.T, fs *utils.InMemoryFS) {
	t.Helper()
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewayCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewayDeploySpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewaySvcSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
}
