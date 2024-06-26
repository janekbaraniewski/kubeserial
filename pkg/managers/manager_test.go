package managers

import (
	"context"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	req := &v1alpha1.ManagerScheduleRequest{}
	manager := &v1alpha1.Manager{
		Spec: v1alpha1.ManagerSpec{
			Config: "dummy",
		},
	}
	api := kubeapi.NewFakeAPIClient()

	err := Schedule(context.TODO(), fs, req, manager, "kubeserial", api)

	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	manager := Manager{}

	err := manager.Delete(context.TODO(), &v1alpha1.KubeSerial{}, &v1alpha1.SerialDevice2{}, kubeapi.NewFakeAPIClient())

	assert.NoError(t, err)
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
