package managers

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controller/api"
)

type Manager struct {
	Image 	string
	RunCmnd	string
}


func (m *Manager) Schedule(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	cm := m.CreateConfigMap(cr, device)
	deploy := m.CreateDeployment(cr, device)
	svc := m.CreateService(cr, device)

	if err := api.EnsureConfigMap(cr, cm); err != nil {
		return err
	}

	if err := api.EnsureDeployment(cr, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(cr, svc); err != nil {
		return err
	}

	if cr.Spec.Ingress.Enabled {
		ingress := m.CreateIngress(cr, device, cr.Spec.Ingress.Domain)
		if err := api.EnsureIngress(cr, ingress); err != nil {
			return err
		}
	}

	return nil
}


func (m *Manager) Delete(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	name := m.GetName(cr.Name, device.Name)

	if err := api.DeleteDeployment(cr, name); err != nil {
		return err
	}
	if err := api.DeleteConfigMap(cr, name); err != nil {
		return err
	}
	if err := api.DeleteService(cr, name); err != nil {
		return err
	}
	if err := api.DeleteIngress(cr, name); err != nil {
		return err
	}

	return nil
}
