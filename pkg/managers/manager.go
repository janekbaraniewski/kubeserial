package managers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	api "github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/janekbaraniewski/kubeserial/pkg/utils/apis"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Manager struct {
	Image      string
	RunCmnd    string
	Config     string
	ConfigPath string
	FS         utils.FileSystem
}

func Schedule(
	ctx context.Context,
	fs utils.FileSystem,
	request *appv1alpha1.ManagerScheduleRequest,
	mgr *appv1alpha1.Manager,
	namespace string,
	api api.API,
) error {
	manager := &Manager{
		Image:      mgr.Spec.Image.Repository + ":" + mgr.Spec.Image.Tag,
		RunCmnd:    mgr.Spec.RunCmd,
		Config:     mgr.Spec.Config,
		ConfigPath: mgr.Spec.ConfigPath,
		FS:         fs,
	}
	cr := types.NamespacedName{
		Name:      request.Name,
		Namespace: namespace,
	}
	if mgr.Spec.Config != "" {
		cm, err := manager.CreateConfigMap(cr, request.Spec.Device)
		if err != nil {
			return err
		}
		if err := api.EnsureObject(ctx, request, cm); err != nil {
			return err
		}
	}
	deploy, err := manager.CreateDeployment(cr, request.Spec.Device, mgr.Spec.Config != "")
	if err != nil {
		return err
	}

	svc, err := manager.CreateService(cr, request.Spec.Device)
	if err != nil {
		return err
	}

	if err := api.EnsureObject(ctx, request, deploy); err != nil {
		return err
	}

	if err := api.EnsureObject(ctx, request, svc); err != nil {
		return err
	}

	return nil
}

func (m *Manager) CreateConfigMap(cr types.NamespacedName, deviceName string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	if err := utils.LoadResourceFromYaml(m.FS, kubeserial.ManagerCMSpecPath, cm); err != nil {
		return cm, err
	}
	name := m.GetName(cr.Name, deviceName)

	cm.Labels[string(kubeserial.AppNameLabel)] = name
	cm.Name = name

	cm.Data = map[string]string{
		filepath.Base(m.ConfigPath): m.Config,
	}

	return cm, nil
}

func (m *Manager) CreateDeployment(cr types.NamespacedName, deviceName string, includeCM bool) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	if err := utils.LoadResourceFromYaml(m.FS, kubeserial.ManagerDeploySpecPath, deployment); err != nil {
		return deployment, err
	}
	name := m.GetName(cr.Name, deviceName)
	deployment.Name = name
	deployment.Labels[string(kubeserial.AppNameLabel)] = name
	deployment.Spec.Selector.MatchLabels[string(kubeserial.AppNameLabel)] = name
	deployment.Spec.Template.Labels[string(kubeserial.AppNameLabel)] = name
	deployment.Spec.Template.Name = name

	deployment.Spec.Template.Spec.Containers[0].Image = m.Image
	deployment.Spec.Template.Spec.Containers[0].Args = []string{
		"-c",
		fmt.Sprintf(
			"socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:%v:3333 & %v",
			strings.ToLower(apis.GatewayName(deviceName)), m.RunCmnd,
		),
	}

	if !includeCM {
		return deployment, nil
	}

	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, corev1.Volume{
		Name: "config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
				Items: []corev1.KeyToPath{
					{
						Key:  filepath.Base(m.ConfigPath),
						Path: filepath.Base(m.ConfigPath),
					},
				},
			},
		},
	})

	container := deployment.Spec.Template.Spec.Containers[0]
	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "config",
			ReadOnly:  false,
			MountPath: m.ConfigPath,
			SubPath:   filepath.Base(m.ConfigPath),
		},
	}

	deployment.Spec.Template.Spec.Containers = []corev1.Container{container}

	return deployment, nil
}

func (m *Manager) CreateService(cr types.NamespacedName, deviceName string) (*corev1.Service, error) {
	svc := &corev1.Service{}
	if err := utils.LoadResourceFromYaml(m.FS, kubeserial.ManagerSvcSpecPath, svc); err != nil {
		return svc, err
	}

	name := m.GetName(cr.Name, deviceName)
	svc.Name = name
	svc.Labels[string(kubeserial.AppNameLabel)] = name
	svc.Spec.Selector[string(kubeserial.AppNameLabel)] = name

	return svc, nil
}
