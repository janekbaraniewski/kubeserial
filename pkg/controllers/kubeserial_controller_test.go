package controllers

import (
	"context"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var ConfigMapGVK = schema.GroupVersionKind{
	Group:   "",
	Version: "v1",
	Kind:    "ConfigMap",
}

var DaemonSetGVK = schema.GroupVersionKind{
	Group:   "apps",
	Version: "v1",
	Kind:    "DaemonSet",
}

func getCR() *kubeserialv1alpha1.KubeSerial {
	return &kubeserialv1alpha1.KubeSerial{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubeserialtest",
			Namespace: "kubeserial",
		},
		Spec: kubeserialv1alpha1.KubeSerialSpec{
			SerialDevices: []kubeserialv1alpha1.SerialDevice2{
				{
					Name:      "testDevice",
					IDVendor:  "0",
					IDProduct: "1",
					Manager:   "testManager",
				},
			},
		},
	}
}

func TestReconcile(t *testing.T) {
	// t.Parallel()

	{
		t.Run("object-not-found", func(t *testing.T) {
			// t.Parallel()

			scheme := runtime.NewScheme()
			utilruntime.Must(clientgoscheme.AddToScheme(scheme))
			utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
			fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

			reconciler := KubeSerialReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeAPIClient(),
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "kubeserial",
					Name:      "kubeserialtest",
				},
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcileReq)
			require.NoError(t, err)
			assert.Equal(t, reconcile.Result{}, result)
		})
	}
	{
		t.Run("object-found", func(t *testing.T) {
			// t.Parallel()

			scheme := runtime.NewScheme()
			utilruntime.Must(clientgoscheme.AddToScheme(scheme))
			utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
			fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
			cr := getCR()
			fs := GetFileSystem(t)

			reconciler := KubeSerialReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				FS:        fs,
				APIClient: kubeapi.NewFakeAPIClient(),
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "kubeserial",
					Name:      "kubeserialtest",
				},
			}
			//nolint:errcheck
			fakeClient.Create(context.TODO(), cr)
			result, err := reconciler.Reconcile(context.TODO(), reconcileReq)
			require.NoError(t, err)
			assert.Equal(t, reconcile.Result{}, result)
		})
	}
}

func GetFileSystem(t *testing.T) utils.FileSystem {
	t.Helper()
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorDSSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	return fs
}

func TestReconcileMonitor(t *testing.T) {
	// t.Parallel()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	fs := GetFileSystem(t)
	apiClient := kubeapi.NewFakeAPIClient()
	reconciler := KubeSerialReconciler{
		Client:    fakeClient,
		Scheme:    scheme,
		FS:        fs,
		APIClient: apiClient,
	}

	err := reconciler.ReconcileMonitor(context.TODO(), getCR(), "latest")

	expected := []kubeapi.Operation{
		{
			Action: "EnsureObject",
			GVK:    ConfigMapGVK,
		},
		{
			Action: "EnsureObject",
			GVK:    DaemonSetGVK,
		},
	}

	require.NoError(t, err)
	assert.Equal(t, expected, apiClient.Operations)
}
