package gateway

import (
	"fmt"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(device metav1.Object, fs utils.FileSystem) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	name := fmt.Sprintf("%v-gateway", device.GetName())

	conf := fmt.Sprintf("3333:raw:600:/dev/%v:115200 8DATABITS NONE 1STOPBIT -XONXOFF LOCAL -RTSCTS HANGUP_WHEN_DONE\n", device.GetName())

	if err := utils.LoadResourceFromYaml(fs, kubeserial.GatewayCMSpecPath, cm); err != nil {
		return cm, err
	}

	cm.ObjectMeta.Labels["app.kubernetes.io/name"] = name
	cm.ObjectMeta.Name = name
	cm.Data["ser2net.conf"] = conf

	return cm, nil
}

func CreateDeployment(device *appv1alpha1.SerialDevice, namespace string, fs utils.FileSystem) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}

	if err := utils.LoadResourceFromYaml(fs, kubeserial.GatewayDeploySpecPath, deployment); err != nil {
		return deployment, err
	}
	name := fmt.Sprintf("%v-gateway", device.GetName())

	deployment.ObjectMeta.Name = name
	deployment.ObjectMeta.Labels["app.kubernetes.io/name"] = name
	deployment.Spec.Selector.MatchLabels["app.kubernetes.io/name"] = name
	deployment.Spec.Template.ObjectMeta.Labels["app.kubernetes.io/name"] = name
	deployment.Spec.Template.ObjectMeta.Name = name

	deployment.Spec.Template.Spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname": device.Status.NodeName,
	}
	volumes := []corev1.Volume{}

	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "config" {
			volume.ConfigMap.Name = name
		}
		volumes = append(volumes, volume)
	}
	deployment.Spec.Template.Spec.Volumes = volumes

	return deployment, nil
}

func CreateService(device *appv1alpha1.SerialDevice, namespace string, fs utils.FileSystem) (*corev1.Service, error) {
	svc := &corev1.Service{}
	if err := utils.LoadResourceFromYaml(fs, kubeserial.GatewaySvcSpecPath, svc); err != nil {
		return svc, err
	}
	name := fmt.Sprintf("%v-gateway", device.GetName())
	svc.ObjectMeta.Name = name
	svc.ObjectMeta.Labels["app.kubernetes.io/name"] = name
	svc.Spec.Selector["app.kubernetes.io/name"] = name
	return svc, nil
}
