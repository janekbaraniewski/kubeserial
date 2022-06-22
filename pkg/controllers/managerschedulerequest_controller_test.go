package controllers

import (
	"context"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	controllerruntime "sigs.k8s.io/controller-runtime"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestManagerScheduleRequestReconcile(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

	reconciler := ManagerScheduleRequestReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	reconciler.Reconcile(context.TODO(), controllerruntime.Request{})
}

func TestManagerScheduleRequestReconcile_DeviceFoundNoManager(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	msr := &v1alpha1.ManagerScheduleRequest{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-schedule-request",
		},
		Spec: v1alpha1.ManagerScheduleRequestSpec{
			Device: "test-device",
		},
	}

	device := &v1alpha1.SerialDevice{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-device",
		},
	}

	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	fakeClient.Create(context.TODO(), msr)
	fakeClient.Create(context.TODO(), device)

	reconciler := ManagerScheduleRequestReconciler{
		Client: fakeClient,
		Scheme: scheme,
	}

	reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: types.NamespacedName{
		Name: "test-schedule-request",
	}})
}

func TestManagerScheduleRequestReconcile_DeviceManagerFound(t *testing.T) {
	t.Parallel()
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	msr := &v1alpha1.ManagerScheduleRequest{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-schedule-request",
		},
		Spec: v1alpha1.ManagerScheduleRequestSpec{
			Device:  "test-device",
			Manager: "test-manager",
		},
	}

	device := &v1alpha1.SerialDevice{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-device",
		},
	}

	manager := &v1alpha1.Manager{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-manager",
		},
	}

	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	fakeClient.Create(context.TODO(), msr)
	fakeClient.Create(context.TODO(), device)
	fakeClient.Create(context.TODO(), manager)

	fs := utils.NewInMemoryFS()
	AddSpecFilesToFilesystem(t, fs)

	apiClient := kubeapi.NewFakeAPIClient()

	reconciler := ManagerScheduleRequestReconciler{
		Client:    fakeClient,
		Scheme:    scheme,
		FS:        fs,
		APIClient: apiClient,
	}

	reconciler.Reconcile(context.TODO(), controllerruntime.Request{NamespacedName: types.NamespacedName{
		Name: "test-schedule-request",
	}})
}

func AddSpecFilesToFilesystem(t *testing.T, fs *utils.InMemoryFS) {
	t.Helper()
	if err := fs.AddFileFromHostPath(string(kubeserial.ManagerCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.ManagerDeploySpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.ManagerSvcSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
}
