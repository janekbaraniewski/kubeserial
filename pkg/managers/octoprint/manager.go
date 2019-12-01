package octoprint

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controller/api"
)

func Schedule(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	cm := CreateConfigMap(cr, device)
	deploy := CreateDeployment(cr, device)
	svc := CreateService(cr, device)

	if err := api.EnsureConfigMap(cr, cm); err != nil {
		return err
	}

	if err := api.EnsureDeployment(cr, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(cr, svc); err != nil {
		return err
	}

	return nil
}

func Delete(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	name := cr.Name + "-" + device.Name + "-manager"  // TODO: this should be set level above (1 place for all types of managers)

	if err := api.DeleteDeployment(cr, name); err != nil {
		return err
	}
	if err := api.DeleteConfigMap(cr, name); err != nil {
		return err
	}
	if err := api.DeleteService(cr, name); err != nil {
		return err
	}

	return nil
}
