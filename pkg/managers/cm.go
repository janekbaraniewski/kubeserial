package managers

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (m *Manager) CreateConfigMap(cr types.NamespacedName, device types.NamespacedName) *corev1.ConfigMap {
	labels := map[string]string{
		"app": m.GetName(cr.Name, device.Name),
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.GetName(cr.Name, device.Name),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Data: map[string]string{
			filepath.Base(m.ConfigPath): m.Config,
		},
	}
}
