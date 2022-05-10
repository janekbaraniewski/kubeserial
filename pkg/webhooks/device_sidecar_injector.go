package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/janekbaraniewski/kubeserial/pkg/injector"
)

// +kubebuilder:webhook:path=/mutate-add-sidecar,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=injector.kubeserial.com,admissionReviewVersions={v1, v1beta1},sideEffects=None

var log = logf.Log.WithName("DeviceSidecarInjecttor")

const (
	requestDeviceSidecarAnnotation   = "app.kubeserial.com/inject-device"
	sidecarAlreadyInjectedAnnotation = "app.kubeserial.com/device-injected"
)

type SidecarInjector struct {
	Name                string
	Client              client.Client
	decoder             *admission.Decoder
	KubeSerialNamespace string
}

func shoudInject(pod *corev1.Pod) string {
	deviceToInject := pod.Annotations[requestDeviceSidecarAnnotation]

	if deviceToInject == "" {
		return deviceToInject
	}

	alreadyUpdated, err := strconv.ParseBool(pod.Annotations[sidecarAlreadyInjectedAnnotation])

	if err == nil && alreadyUpdated {
		return ""
	}

	log.Info("Should Inject", "device to inject", deviceToInject)

	return deviceToInject
}

// SidecarInjector adds an annotation to every incoming pods.
func (si *SidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := si.decoder.Decode(req, pod)
	if err != nil {
		log.Info("Sdecar-Injector: cannot decode")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	deviceToInject := shoudInject(pod)

	if deviceToInject != "" {
		deviceGateway := types.NamespacedName{
			Name:      fmt.Sprintf("%v-gateway", deviceToInject),
			Namespace: si.KubeSerialNamespace,
		}
		log.Info("Injecting device mounter connecting to gateway", "gateway", deviceGateway)
		injector.AddDeviceInjector(&pod.Spec, deviceGateway)
		pod.Annotations[sidecarAlreadyInjectedAnnotation] = "true"
		log.Info("Sidecar injected", "sidecar name", si.Name)
	} else {
		log.Info("Inject not needed.")
	}

	marshaledPod, err := json.Marshal(pod)

	if err != nil {
		log.Info("Sdecar-Injector: cannot marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// SidecarInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *SidecarInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
