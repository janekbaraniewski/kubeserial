package controllers

import (
	"context"
	"testing"
	"time"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
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
			SerialDevices: []kubeserialv1alpha1.SerialDevice_2{
				{
					Name:      "testDevice",
					IdVendor:  "0",
					IdProduct: "1",
					Manager:   "octoprint", // TODO: this should use some test manager
					Subsystem: "tty",
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
				Client: fakeClient,
				Scheme: scheme,
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "test-namespace",
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
			reconciler := KubeSerialReconciler{
				Client: fakeClient,
				Scheme: scheme,
			}

			reconcileReq := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "test-namespace",
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

func TestReconcileMonitor(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

	reconciler := KubeSerialReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}
	apiClient := api.NewFakeApiClient()

	err := reconciler.ReconcileMonitor(context.TODO(), getCR(), &apiClient, "latest")

	assert.Equal(t, nil, err)
	assert.Equal(t, []string{"EnsureConfigMap", "EnsureDaemonSet"}, apiClient.Operations)
}
