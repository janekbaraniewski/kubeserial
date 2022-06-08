package controllers

import (
	"context"
	"testing"
	"time"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
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
			SerialDevices: []kubeserialv1alpha1.SerialDevice_2{
				{
					Name:      "testDevice",
					IdVendor:  "0",
					IdProduct: "1",
					Manager:   "testManager",
				},
			},
		},
	}
}

func TestReconcile(t *testing.T) {
	{
		t.Run("object-not-found", func(t *testing.T) {
			scheme := runtime.NewScheme()
			utilruntime.Must(clientgoscheme.AddToScheme(scheme))
			utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
			fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

			reconciler := KubeSerialReconciler{
				Client:    fakeClient,
				Scheme:    scheme,
				APIClient: kubeapi.NewFakeApiClient(),
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "kubeserial",
					Name:      "kubeserialtest",
				},
			}

			result, err := reconciler.Reconcile(context.TODO(), reconcileReq)
			assert.Equal(t, nil, err)
			assert.Equal(t, reconcile.Result{}, result)
		})
	}
	{
		t.Run("object-found", func(t *testing.T) {
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
				APIClient: kubeapi.NewFakeApiClient(),
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "kubeserial",
					Name:      "kubeserialtest",
				},
			}
			fakeClient.Create(context.TODO(), cr)
			result, err := reconciler.Reconcile(context.TODO(), reconcileReq)
			assert.Equal(t, nil, err)
			assert.Equal(t, reconcile.Result{RequeueAfter: time.Second * 5}, result)
		})
	}
}

func GetFileSystem(t *testing.T) utils.FileSystem {
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath("monitor-daemonset.yaml"); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath("monitor-configmap.yaml"); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	return fs
}

func TestReconcileMonitor(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	fs := GetFileSystem(t)
	apiClient := api.NewFakeApiClient()
	reconciler := KubeSerialReconciler{
		Client:    fakeClient,
		Scheme:    scheme,
		FS:        fs,
		APIClient: apiClient,
	}

	err := reconciler.ReconcileMonitor(context.TODO(), getCR(), "latest")

	expected := []api.Operation{
		{
			Action: "EnsureObject",
			GVK:    ConfigMapGVK,
		},
		{
			Action: "EnsureObject",
			GVK:    DaemonSetGVK,
		},
	}

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, apiClient.Operations)
}
