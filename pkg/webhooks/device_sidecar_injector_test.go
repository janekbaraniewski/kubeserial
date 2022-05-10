package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	"gomodules.xyz/jsonpatch/v2"
	v1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func TestHandle(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()

	si := SidecarInjector{
		Name:                "TestDeviceSidecarInjector",
		Client:              fakeClient,
		KubeSerialNamespace: "test-ns",
	}

	decoder, _ := admission.NewDecoder(scheme)

	si.InjectDecoder(decoder)

	testPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Annotations: map[string]string{
				requestDeviceSidecarAnnotation: "test-device",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "test-container",
				},
			},
		},
	}

	rawTestPod, _ := json.Marshal(testPod)

	req := admission.Request{
		AdmissionRequest: v1.AdmissionRequest{
			UID: types.UID(""),
			Kind: metav1.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Pod",
			},
			Resource: metav1.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "pods",
			},
			Operation: "CREATE",
			UserInfo:  authenticationv1.UserInfo{},
			Object: runtime.RawExtension{
				Raw:    rawTestPod,
				Object: &testPod,
			},
		},
	}

	resp := si.Handle(context.TODO(), req)

	resultPod := &corev1.Pod{}

	json.Unmarshal(resp.Patch, resultPod)

	assert.Equal(t, true, resp.Allowed)
	assert.Equal(t, true, reflect.DeepEqual(resp.Patches, patches()))
}

func patches() []jsonpatch.Operation {
	return []jsonpatch.Operation{
		{
			Operation: "add",
			Path:      fmt.Sprintf("/metadata/annotations/%v", strings.Replace(sidecarAlreadyInjectedAnnotation, "/", "~1", -1)),
			Value:     "true",
		},
		{

			Operation: "add",
			Path:      "/spec/volumes",
			Value: []interface{}{
				map[string]interface{}{
					"name":     "devices",
					"emptyDir": map[string]interface{}{},
				},
			},
		},
		{
			Operation: "add",
			Path:      "/spec/containers/1",
			Value: map[string]interface{}{
				"args": []interface{}{
					"-c",
					"sleep 5 && socat -d -d pty,raw,echo=0,b115200,link=/dev/devices/test-device-gateway,perm=0660,group=tty tcp:test-device-gateway.test-ns:3333",
				},
				"command": []interface{}{
					"/bin/sh",
				},
				"image":     "alpine/socat:1.7.4.3-r0",
				"name":      "device-mounter",
				"resources": map[string]interface{}{},
				"volumeMounts": []interface{}{
					map[string]interface{}{
						"mountPath": "/dev/devices",
						"name":      "devices",
					},
				},
			},
		},
		{
			Operation: "add",
			Path:      "/spec/containers/0/volumeMounts",
			Value: []interface{}{
				map[string]interface{}{
					"name":      "devices",
					"mountPath": "/dev/devices",
				},
			},
		},
	}
}
