package controllers

import (
	"context"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
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
	deviceName := types.NamespacedName{
		Name:      "test-device",
		Namespace: "test-ns",
	}

	device := &v1alpha1.Device{
		ObjectMeta: v1.ObjectMeta{
			Name:      deviceName.Name,
			Namespace: deviceName.Namespace,
		},
		Spec: v1alpha1.DeviceSpec{
			Manager:   "test-manager",
			IdVendor:  "123",
			IdProduct: "456",
			Subsystem: "tty",
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

	{
		t.Run("device-new-manager-not-available", func(t *testing.T) {
			fakeClient.Create(context.TODO(), device)

			deviceReconciler := DeviceReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.Device{}
			err = fakeClient.Get(context.TODO(), deviceName, foundDevice)

			assert.Equal(t, nil, err)
			assert.Equal(t, 2, len(foundDevice.Status.Conditions))

			availableCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceAvailable)
			assert.Equal(t, v1.ConditionUnknown, availableCondition.Status)
			assert.Equal(t, "NotValidated", availableCondition.Reason)

			readyCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceReady)
			assert.Equal(t, v1.ConditionFalse, readyCondition.Status)
			assert.Equal(t, "ManagerNotAvailable", readyCondition.Reason)
		})
	}
	{
		t.Run("device-new-manager-available", func(t *testing.T) {
			fakeClient.Create(context.TODO(), device)
			fakeClient.Create(context.TODO(), manager)

			deviceReconciler := DeviceReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.Device{}
			fakeClient.Get(context.TODO(), deviceName, foundDevice)

			availableCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceAvailable)
			assert.Equal(t, v1.ConditionUnknown, availableCondition.Status)
			assert.Equal(t, "NotValidated", availableCondition.Reason)

			readyCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceReady)
			assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
			assert.Equal(t, "AllChecksPassed", readyCondition.Reason)
		})
	}
	{
		t.Run("device-available-manager-available", func(t *testing.T) {
			device.Status.Conditions = append(device.Status.Conditions, v1alpha1.DeviceCondition{
				Type:   v1alpha1.DeviceAvailable,
				Status: v1.ConditionTrue,
			})
			fakeClient.Create(context.TODO(), device)
			fakeClient.Create(context.TODO(), manager)

			deviceReconciler := DeviceReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			result, err := deviceReconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: deviceName})

			assert.Equal(t, nil, err)
			assert.Equal(t, false, result.Requeue)

			foundDevice := &v1alpha1.Device{}
			fakeClient.Get(context.TODO(), deviceName, foundDevice)

			availableCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceAvailable)
			assert.Equal(t, v1.ConditionUnknown, availableCondition.Status)
			assert.Equal(t, "NotValidated", availableCondition.Reason)

			readyCondition := utils.GetCondition(foundDevice.Status.Conditions, v1alpha1.DeviceReady)
			assert.Equal(t, v1.ConditionTrue, readyCondition.Status)
			assert.Equal(t, "AllChecksPassed", readyCondition.Reason)

			// TODO: Get schedule request and assert
		})
	}
}
