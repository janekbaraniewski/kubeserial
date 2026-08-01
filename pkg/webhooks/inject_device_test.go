package webhooks

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/fake"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.Install(scheme))
}

func podWithAnnotations(annotations map[string]string) *corev1.Pod {
	return &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Annotations: annotations}}
}

func TestShouldInject(t *testing.T) {
	tests := []struct {
		name string
		pod  *corev1.Pod
		want string
	}{
		{
			name: "device requested",
			pod: podWithAnnotations(map[string]string{
				requestDeviceSidecarAnnotation: "test-device",
			}),
			want: "test-device",
		},
		{
			name: "already injected",
			pod: podWithAnnotations(map[string]string{
				requestDeviceSidecarAnnotation:   "test-device",
				sidecarAlreadyInjectedAnnotation: "true",
			}),
			want: "",
		},
		{
			name: "already-injected marker is not a bool",
			pod: podWithAnnotations(map[string]string{
				requestDeviceSidecarAnnotation:   "test-device",
				sidecarAlreadyInjectedAnnotation: "not-a-bool",
			}),
			want: "test-device",
		},
		{
			name: "no annotations",
			pod:  podWithAnnotations(map[string]string{}),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, shoudInject(tt.pod))
		})
	}
}

func TestDevicePath(t *testing.T) {
	assert.Equal(t, defaultDevicePath, devicePath(podWithAnnotations(nil)),
		"no annotation should fall back to the default path")
	assert.Equal(t, "/dev/ttyMY0", devicePath(podWithAnnotations(map[string]string{
		devicePathAnnotation: "/dev/ttyMY0",
	})), "annotation should override the default path")
}

func TestConcatCommandWithSocat(t *testing.T) {
	cmd, args := concatCommandWithSocat(
		[]string{"testCommand"}, []string{"test", "args"}, "test-device", "/dev/device")

	assert.Equal(t, []string{"/bin/sh"}, cmd)
	require.Len(t, args, 2)
	assert.Equal(t, "-c", args[0])
	assert.Equal(t,
		"socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty "+
			"tcp:test-device-gateway:3333 & testCommand test args",
		args[1])
}

func TestConcatCommandWithSocatHonoursPath(t *testing.T) {
	_, args := concatCommandWithSocat([]string{"run"}, nil, "dev1", "/dev/custom")
	assert.Contains(t, args[1], "link=/dev/custom")
	assert.NotContains(t, args[1], "link=/dev/device")
}

// injectorFor builds an injector backed by a fake kubeserial clientset holding
// the given devices. ConfigExtractor is left nil: every test pod sets an
// explicit command, so the registry-lookup path is never taken.
func injectorFor(devices ...*kubeserialv1alpha1.SerialDevice) *SerialDeviceInjector {
	objs := make([]runtime.Object, 0, len(devices))
	for _, d := range devices {
		objs = append(objs, d)
	}
	return &SerialDeviceInjector{
		Name:      "test-injector",
		Clientset: fake.NewSimpleClientset(objs...),
		Decoder:   admission.NewDecoder(scheme),
	}
}

func deviceWithFree(name string, status metav1.ConditionStatus) *kubeserialv1alpha1.SerialDevice {
	dev := &kubeserialv1alpha1.SerialDevice{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: kubeserialv1alpha1.SerialDeviceSpec{
			Name: name, IDVendor: "0403", IDProduct: "6001",
		},
	}
	dev.SetCondition(kubeserialv1alpha1.SerialDeviceCondition{
		Type:   kubeserialv1alpha1.SerialDeviceFree,
		Status: status,
		Reason: "Test",
	})
	return dev
}

func requestFor(t *testing.T, pod *corev1.Pod) admission.Request {
	t.Helper()
	raw, err := json.Marshal(pod)
	require.NoError(t, err)
	return admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Object: runtime.RawExtension{Raw: raw},
		},
	}
}

// applyPatches runs the admission response's JSON patch against the original
// pod so assertions are made on the pod the API server would actually persist.
func applyPatches(t *testing.T, pod *corev1.Pod, resp admission.Response) *corev1.Pod {
	t.Helper()
	require.True(t, resp.Allowed, "response should be allowed")

	original, err := json.Marshal(pod)
	require.NoError(t, err)

	patchJSON, err := json.Marshal(resp.Patches)
	require.NoError(t, err)

	patch, err := jsonpatch.DecodePatch(patchJSON)
	require.NoError(t, err)

	patched, err := patch.Apply(original)
	require.NoError(t, err)

	out := &corev1.Pod{}
	require.NoError(t, json.Unmarshal(patched, out))
	return out
}

func podRequestingDevice(name, device string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "default",
			Annotations: map[string]string{requestDeviceSidecarAnnotation: device},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:    "app",
				Image:   "busybox:latest",
				Command: []string{"/bin/sh"},
				Args:    []string{"-c", "sleep 3600"},
			}},
		},
	}
}

func containerCmdline(pod *corev1.Pod) string {
	c := pod.Spec.Containers[0]
	return strings.Join(append(c.Command, c.Args...), " ")
}

func TestHandleInjectsWhenDeviceFree(t *testing.T) {
	pod := podRequestingDevice("app", "test-device")
	si := injectorFor(deviceWithFree("test-device", metav1.ConditionTrue))

	got := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))

	cmdline := containerCmdline(got)
	assert.Contains(t, cmdline, "socat")
	assert.Contains(t, cmdline, "tcp:test-device-gateway:3333")
	assert.Contains(t, cmdline, "sleep 3600", "original command must be preserved")
	assert.Equal(t, "true", got.Annotations[sidecarAlreadyInjectedAnnotation],
		"injection must mark the pod so UPDATE does not re-wrap it")
}

func TestHandleDoesNotInjectWhenDeviceNotFree(t *testing.T) {
	pod := podRequestingDevice("app", "busy-device")
	si := injectorFor(deviceWithFree("busy-device", metav1.ConditionFalse))

	got := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))

	assert.NotContains(t, containerCmdline(got), "socat")
	assert.NotContains(t, got.Annotations, sidecarAlreadyInjectedAnnotation)
}

func TestHandleDoesNotInjectWhenDeviceMissing(t *testing.T) {
	pod := podRequestingDevice("app", "ghost-device")
	si := injectorFor()

	got := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))

	assert.NotContains(t, containerCmdline(got), "socat")
}

func TestHandleIgnoresPodWithoutAnnotation(t *testing.T) {
	pod := podRequestingDevice("app", "test-device")
	delete(pod.Annotations, requestDeviceSidecarAnnotation)
	si := injectorFor(deviceWithFree("test-device", metav1.ConditionTrue))

	got := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))

	assert.NotContains(t, containerCmdline(got), "socat")
	assert.Contains(t, containerCmdline(got), "sleep 3600")
}

// Regression test for the double-injection bug: an UPDATE of an already
// injected pod must be left alone rather than nesting a second socat bridge.
func TestHandleIsIdempotentAcrossUpdates(t *testing.T) {
	pod := podRequestingDevice("app", "test-device")
	si := injectorFor(deviceWithFree("test-device", metav1.ConditionTrue))

	first := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))
	require.Equal(t, 1, strings.Count(containerCmdline(first), "socat"))

	second := applyPatches(t, first, si.Handle(context.Background(), requestFor(t, first)))

	assert.Equal(t, 1, strings.Count(containerCmdline(second), "socat"),
		"re-admitting an injected pod must not wrap the command a second time")
	assert.Equal(t, containerCmdline(first), containerCmdline(second))
}

func TestHandleHonoursDevicePathAnnotation(t *testing.T) {
	pod := podRequestingDevice("app", "test-device")
	pod.Annotations[devicePathAnnotation] = "/dev/ttyPRINTER"
	si := injectorFor(deviceWithFree("test-device", metav1.ConditionTrue))

	got := applyPatches(t, pod, si.Handle(context.Background(), requestFor(t, pod)))

	assert.Contains(t, containerCmdline(got), "link=/dev/ttyPRINTER")
}

func TestHandleRejectsUndecodableRequest(t *testing.T) {
	si := injectorFor()
	resp := si.Handle(context.Background(), admission.Request{})

	assert.False(t, resp.Allowed)
	require.NotNil(t, resp.Result)
	assert.EqualValues(t, 400, resp.Result.Code)
}
