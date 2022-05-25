package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/janekbaraniewski/kubeserial/pkg/images"
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
	Name            string
	Client          client.Client
	ConfigExtractor *images.OCIConfigExtractor
	decoder         *admission.Decoder
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

	if deviceToInject == "" {
		log.Info("Inject not needed.")

		marshaledPod, err := json.Marshal(pod)

		if err != nil {
			log.Info("Sdecar-Injector: cannot marshal")
			return admission.Errored(http.StatusInternalServerError, err)
		}

		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
	}

	log.Info("Pod is requesting device, checking if available", "pod", pod.Name, "device", deviceToInject)
	log.Info("FAKE device available") // TODO: check if device available, for now happy path

	container := &pod.Spec.Containers[0] // TODO: what if there are multiple containers? probably should introduce some annotation to select one

	command, args := getContainerCommandArgs(container)
	log.Info(
		"Manager pod command and args",
		"command", command,
		"args", args,
	)

	if command == nil {
		// TODO: implement overriding image entrypoint
		log.Info("Image", "image", container.Image)

		imageConfig, err := si.ConfigExtractor.GetImageConfig(ctx, container.Image)
		if err != nil {
			panic(err)
		}

		log.Info(
			"Manager container entrypoint and cmd",
			"entrypoint", imageConfig.Entrypoint,
			"cmd", imageConfig.Cmd,
		)

		command = imageConfig.Entrypoint
		args = imageConfig.Cmd
	}

	newCommand, newArgs := concatCommandWithSocat(command, args, deviceToInject)

	container.Command = newCommand
	container.Args = newArgs
	//TODO: mutate command and args, maybe the best would be to mount entrypoint from some CM?
	log.Info(
		"Injected",
		"Container command", pod.Spec.Containers[0].Command,
		"Container args", pod.Spec.Containers[0].Args,
	)

	marshaledPod, err := json.Marshal(pod)

	if err != nil {
		log.Info("Sdecar-Injector: cannot marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func concatCommandWithSocat(command []string, args []string, device string) (newCommand []string, newArgs []string) {

	nCommand := []string{"/bin/sh"}
	nArgs := []string{
		"-c",
		fmt.Sprintf("socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:%v-gateway:3333 & %v", device, strings.Join(append(command, args...), " ")),
	}
	return nCommand, nArgs
}

func getContainerCommandArgs(container *corev1.Container) (command []string, args []string) {
	return container.Command, container.Args
}

// DeviceInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *DeviceInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
