package gateway

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateConfigMap(device metav1.Object) *corev1.ConfigMap {
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
