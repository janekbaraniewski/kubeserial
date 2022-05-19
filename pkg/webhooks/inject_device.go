package webhooks

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-mount-device,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=device.kubeserial.com,admissionReviewVersions={v1, v1beta1},sideEffects=None

var log = logf.Log.WithName("DeviceSidecarInjecttor")

const (
	requestDeviceSidecarAnnotation   = "app.kubeserial.com/inject-device"
	sidecarAlreadyInjectedAnnotation = "app.kubeserial.com/device-injected"
)

type DeviceInjector struct {
	Name    string
	Client  client.Client
	decoder *admission.Decoder
	Config  *Config
}

type Config struct {
	Containers  []corev1.Container `yaml:"containers"`
	Volume      corev1.Volume      `yaml:"volume"`
	VolumeMount corev1.VolumeMount `yaml:"volumeMount"`
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

// DeviceInjector mutates command and args to inject script that mounts selected device.
// It checks if pod requested device and if requested device is available.
func (si *DeviceInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
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
		log.Info("Injecting volume...", "volumeConfig", si.Config.Volume)
		pod.Spec.Volumes = append(pod.Spec.Volumes, si.Config.Volume)
		log.Info("Injecting sidecar...", "sidecarConfig", si.Config.Containers)
		pod.Spec.Containers = append(pod.Spec.Containers, si.Config.Containers...)
		containers := []corev1.Container{}
		for _, container := range pod.Spec.Containers {
			container.VolumeMounts = append(container.VolumeMounts, si.Config.VolumeMount)
			log.Info("Attached volume mounts", "volumeMounta", container.VolumeMounts)
			containers = append(containers, container)
		}
		pod.Spec.Containers = containers
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

// DeviceInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *DeviceInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
