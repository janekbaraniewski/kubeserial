package gateway

import (
	"strings"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device_2) *corev1.ConfigMap {
	labels := map[string]string{
		"app": strings.ToLower(cr.Name + "-" + device.Name + "-gateway"),
	}

	conf := "3333:raw:600:/dev/" + device.Name + ":115200 8DATABITS NONE 1STOPBIT -XONXOFF LOCAL -RTSCTS HANGUP_WHEN_DONE\n"

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(cr.Name + "-" + device.Name + "-gateway"),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			"ser2net.conf": conf,
		},
	}
}

func CreateConfigMapNew(device metav1.Object) *corev1.ConfigMap {
	labels := map[string]string{
		"app": strings.ToLower(device.GetName() + "-gateway"),
	}

	conf := "3333:raw:600:/dev/" + device.GetName() + ":115200 8DATABITS NONE 1STOPBIT -XONXOFF LOCAL -RTSCTS HANGUP_WHEN_DONE\n"

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(device.GetName() + "-gateway"),
			Namespace: device.GetNamespace(),
			Labels:    labels,
		},
		Data: map[string]string{
			"ser2net.conf": conf,
		},
	}
}
