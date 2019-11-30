package gateway

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func CreateConfigMap(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) *corev1.ConfigMap {
	labels := map[string]string{
		"app": cr.Name + "-" + device.Name + "-gateway",
	}

	conf := "3333:raw:600:/dev/tty" + device.Name + ":115200 8DATABITS NONE 1STOPBIT -XONXOFF LOCAL -RTSCTS HANGUP_WHEN_DONE\n"

	return &corev1.ConfigMap {
		ObjectMeta:		metav1.ObjectMeta {
			Name: 		cr.Name + "-" + device.Name + "-gateway",
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Data:			map[string]string {
			"ser2net.conf": conf,
		},
	}
}
