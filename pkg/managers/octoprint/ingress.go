package octoprint

import (
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateIngress(cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device) *v1beta1.Ingress {
	labels := map[string]string{
		"app": cr.Name + "-" + device.Name + "-manager",
	}
	return &v1beta1.Ingress{
		ObjectMeta:	metav1.ObjectMeta {
			Name:		cr.Name + "-" + device.Name + "-manager",
			Namespace:	cr.Namespace,
			Labels:		labels,
		},
		Spec:		v1beta1.IngressSpec{
			Rules:		[]v1beta1.IngressRule{
				{
					Host:				device.Name + ".my.home",  // TODO: parametrize
					IngressRuleValue:	v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths:	[]v1beta1.HTTPIngressPath{
								{
									Path: 		"/",
									Backend: 	v1beta1.IngressBackend {
										ServiceName: cr.Name + "-" + device.Name + "-manager",
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
