package webhooks

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/ref"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-inject-device,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=device.kubeserial.com,admissionReviewVersions={v1, v1beta1},sideEffects=None

var log = logf.Log.WithName("DeviceSidecarInjecttor")

const (
	requestDeviceSidecarAnnotation   = "app.kubeserial.com/inject-device"
	sidecarAlreadyInjectedAnnotation = "app.kubeserial.com/device-injected"
)

type DeviceInjector struct {
	Name    string
	Client  client.Client
	decoder *admission.Decoder
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
		log.Info("Pod is requesting device, checking if available", "pod", pod.Name, "device", deviceToInject)
		// TODO: check if device available, for now happy path
		log.Info("FAKE device available")
		log.Info(
			"Manager pod command and args",
			"command", pod.Spec.Containers[0].Command,
			"args", pod.Spec.Containers[0].Args,
		)

		client := regclient.New()
		ref, err := ref.New(pod.Spec.Containers[0].Image)

		if err != nil {
			panic(err)
		}

		manifest, err := client.ManifestGet(ctx, ref)

		if err != nil {
			panic(err)
		}

		config, err := manifest.GetConfig()

		if err != nil {
			panic(err)
		}

		blobConfig, err := client.BlobGetOCIConfig(ctx, ref, config)

		if err != nil {
			panic(err)
		}

		log.Info(
			"Manager container entrypoint and cmd",
			"entrypoint", blobConfig.GetConfig().Config.Entrypoint,
			"cmd", blobConfig.GetConfig().Config.Cmd,
		)

		//TODO: mutate command and args, maybe the best would be to mount entrypoint from some CM?
		log.Info("Injected")
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
