package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/images"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-inject-device,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=device.kubeserial.com,admissionReviewVersions={v1, v1beta1},sideEffects=None

var log = logf.Log.WithName("DeviceSidecarInjecttor")

const (
	requestDeviceSidecarAnnotation   = "app.kubeserial.com/inject-device"
	sidecarAlreadyInjectedAnnotation = "app.kubeserial.com/device-injected"
)

type SerialDeviceInjector struct {
	Name            string
	Clientset       versioned.Interface
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

func returnProperResponse(pod *corev1.Pod, req admission.Request) admission.Response {
	marshaledPod, err := json.Marshal(pod)

	if err != nil {
		log.Info("cannot marshal")
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

// SerialDeviceInjector mutates command and args to inject script that mounts selected device.
// It checks if pod requested device and if requested device is available.
func (si *SerialDeviceInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	PodsHandled.Inc()
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
		return returnProperResponse(pod, req)
	}

	log.Info("Pod is requesting device, checking if available", "pod", pod.Name, "device", deviceToInject)
	log.Info("Looking for device", "device", deviceToInject)

	device, err := si.Clientset.AppV1alpha1().SerialDevices().Get(ctx, deviceToInject, metav1.GetOptions{})

	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "Device not found!", "device", deviceToInject)
		}
		log.Error(err, "Some error when looking for device", "device", deviceToInject)
		return returnProperResponse(pod, req)
	}
	log.Info("Device found, checking if free", "device", device)

	condition := device.GetCondition(v1alpha1.SerialDeviceFree)

	if condition.Status != metav1.ConditionTrue {
		log.Info("Device is not free, not injecting", "device", device, "condition", condition)
		return returnProperResponse(pod, req)
	}

	log.Info("Device is free", "device", device, "condition", condition) // TODO: check if device available, for now happy path

	container := &pod.Spec.Containers[0] // TODO: what if there are multiple containers? probably should introduce some annotation to select one

	command, args := getContainerCommandArgs(container)
	log.Info(
		"Manager pod command and args",
		"command", command,
		"args", args,
	)

	if command == nil {
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
	InjectedCommands.Inc()
	//TODO: mutate command and args, maybe the best would be to mount entrypoint from some CM?
	log.Info(
		"Injected",
		"Container command", pod.Spec.Containers[0].Command,
		"Container args", pod.Spec.Containers[0].Args,
	)

	return returnProperResponse(pod, req)
}

// SerialDeviceInjector implements admission.DecoderInjector.
// A decoder will be automatically inj1ected.

// InjectDecoder injects the decoder.
func (si *SerialDeviceInjector) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
