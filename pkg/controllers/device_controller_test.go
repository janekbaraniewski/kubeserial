package controllers

import (
	"context"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
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
	t.Parallel()
	deviceName := types.NamespacedName{
		Name:      "test-device",
		Namespace: "test-ns",
	}

	device := &v1alpha1.SerialDevice{
		ObjectMeta: v1.ObjectMeta{
			Name:      deviceName.Name,
			Namespace: deviceName.Namespace,
		},
		Spec: v1alpha1.SerialDeviceSpec{
			Manager:   "test-manager",
			IDVendor:  "123",
			IDProduct: "456",
			Name:      "test-device",
		},
	}

	manager := &v1alpha1.Manager{
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

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	fs := utils.NewInMemoryFS()
	AddGatewaySpecFilesToFilesystem(t, fs)
	{
		t.Run("device-new-manager-not-available", func(t *testing.T) {
			t.Parallel()
			t.Skip()
			//nolint:errcheck
			fakeClient.Create(context.TODO(), device)

			deviceReconciler := SerialDeviceReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeAPIClient(),
				FS:        fs,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.SerialDevice{}
			err = fakeClient.Get(context.TODO(), deviceName, foundDevice)

			assert.Equal(t, nil, err)
			assert.Equal(t, 3, len(foundDevice.Status.Conditions))

			availableCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceAvailable)
			assert.Equal(t, v1.ConditionFalse, availableCondition.Status)
			assert.Equal(t, "NotValidated", availableCondition.Reason)

			readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
			assert.Equal(t, v1.ConditionFalse, readyCondition.Status)
			assert.Equal(t, "ManagerNotAvailable", readyCondition.Reason)
		})
	}
	{
		t.Run("device-new-manager-available", func(t *testing.T) {
			t.Parallel()
			t.Skip()
			//nolint:errcheck
			fakeClient.Create(context.TODO(), device)
			//nolint:errcheck
			fakeClient.Create(context.TODO(), manager)

			deviceReconciler := SerialDeviceReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeAPIClient(),
				FS:        fs,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.SerialDevice{}
			//nolint:errcheck
			fakeClient.Get(context.TODO(), deviceName, foundDevice)

			availableCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceAvailable)
			assert.Equal(t, v1.ConditionFalse, availableCondition.Status)
			assert.Equal(t, "NotValidated", availableCondition.Reason)

			readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
			assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
			assert.Equal(t, "AllChecksPassed", readyCondition.Reason)
		})
	}
	{
		t.Run("device-ready", func(t *testing.T) {
			t.Parallel()
			t.Skip()
			device.Status.Conditions = append(device.Status.Conditions, v1alpha1.SerialDeviceCondition{
				Type:   v1alpha1.SerialDeviceAvailable,
				Status: v1.ConditionTrue,
			})
			//nolint:errcheck
			fakeClient.Create(context.TODO(), device)
			//nolint:errcheck
			fakeClient.Create(context.TODO(), manager)

			deviceReconciler := SerialDeviceReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeAPIClient(),
				FS:        fs,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.SerialDevice{}
			//nolint:errcheck
			fakeClient.Get(context.TODO(), deviceName, foundDevice)

			readyCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceReady)
			assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
			assert.Equal(t, "AllChecksPassed", readyCondition.Reason)

			foundRequest := v1alpha1.ManagerScheduleRequest{}
			//nolint:errcheck
			fakeClient.Get(context.TODO(), types.NamespacedName{
				Name:      device.Name + "-" + device.Spec.Manager,
				Namespace: device.Name,
			}, &foundRequest)

			assert.Equal(t, false, foundRequest.Status.Fulfilled)
		})
	}
	{
		t.Run("device-not-found", func(t *testing.T) {
			t.Parallel()

			deviceReconciler := SerialDeviceReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeAPIClient(),
				FS:        fs,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})
			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)
		})
	}
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
