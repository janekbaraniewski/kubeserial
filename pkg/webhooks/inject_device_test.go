package webhooks

// import (
// 	"context"
// 	"errors"
// 	"reflect"
// 	"testing"

// 	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
// 	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
// 	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/fake"
// 	"github.com/janekbaraniewski/kubeserial/pkg/images"
// 	"github.com/stretchr/testify/assert"
// 	admissionv1 "k8s.io/api/admission/v1"
// 	corev1 "k8s.io/api/core/v1"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
// 	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
// 	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
// )

// var (
// 	scheme = runtime.NewScheme()
// )

// func init() {
// 	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

// 	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
// 	//+kubebuilder:scaffold:scheme
// }

// func getPodWithAnnotations(annotations map[string]string) *corev1.Pod {
// 	return &corev1.Pod{
// 		ObjectMeta: v1.ObjectMeta{
// 			Annotations: annotations,
// 		},
// 	}
// }

// func Test_shoudInject(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		pod  *corev1.Pod
// 		want string
// 	}{
// 		{
// 			name: "should_inject_device",
// 			pod: getPodWithAnnotations(map[string]string{
// 				requestDeviceSidecarAnnotation: "test-device",
// 			}),
// 			want: "test-device",
// 		},
// 		{
// 			name: "device_already_injected",
// 			pod: getPodWithAnnotations(map[string]string{
// 				requestDeviceSidecarAnnotation:   "test-device",
// 				sidecarAlreadyInjectedAnnotation: "true",
// 			}),
// 			want: "",
// 		},
// 		{
// 			name: "no_injection_needed",
// 			pod:  getPodWithAnnotations(map[string]string{}),
// 			want: "",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := shoudInject(tt.pod); got != tt.want {
// 				t.Errorf("shoudInject() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_returnProperResponse(t *testing.T) {
// 	type args struct {
// 		pod *corev1.Pod
// 		req admission.Request
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want admission.Response
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := returnProperResponse(tt.args.pod, tt.args.req); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("returnProperResponse() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_concatCommandWithSocat(t *testing.T) {
// 	command := []string{"testCommand"}
// 	args := []string{"test", "args"}
// 	device := "test-device"
// 	newCommand, newArgs := concatCommandWithSocat(command, args, device)
// 	wantNewCommand := []string{"/bin/sh"}
// 	wantNewArgs := []string{
// 		"-c",
// 		"socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:test-device-gateway:3333 & testCommand test args",
// 	}
// 	if !reflect.DeepEqual(newCommand, wantNewCommand) {
// 		t.Errorf("concatCommandWithSocat() gotNewCommand = %v, want %v", newCommand, wantNewCommand)
// 	}
// 	if !reflect.DeepEqual(newArgs, wantNewArgs) {
// 		t.Errorf("concatCommandWithSocat() gotNewArgs = %v, want %v", newArgs, wantNewArgs)
// 	}
// }

// func TestSerialDeviceInjector_Handle(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		req  admission.Request
// 		want admission.Response
// 	}{
// 		{
// 			name: "no_content_to_decode",
// 			req:  admission.Request{},
// 			want: admission.Errored(400, errors.New("there is no content to decode")),
// 		},
// 		{
// 			name: "allowed_when_no_injection_needed",
// 			req: admission.Request{
// 				AdmissionRequest: admissionv1.AdmissionRequest{
// 					Object: runtime.RawExtension{
// 						Raw: []byte(`{
// 			"apiVersion": "v1",
// 			"kind": "Pod",
// 			"metadata": {
// 				"name": "foo",
// 				"namespace": "default"
// 			},
// 			"spec": {
// 				"containers": [
// 					{
// 						"image": "bar:v2",
// 						"name": "bar"
// 					}
// 				]
// 			}
// 		}`),
// 					},
// 				},
// 			},
// 			want: admission.Allowed(""),
// 		},
// 		{
// 			name: "allowed_when_device_not_found",
// 			req: admission.Request{
// 				AdmissionRequest: admissionv1.AdmissionRequest{
// 					Object: runtime.RawExtension{
// 						Raw: []byte(`{
// 			"apiVersion": "v1",
// 			"kind": "Pod",
// 			"metadata": {
// 				"name": "foo",
// 				"namespace": "default",
// 				"annotations": {
// 					"app.kubeserial.com/inject-device": "missing-device"
// 				}
// 			},
// 			"spec": {
// 				"containers": [
// 					{
// 						"image": "bar:v2",
// 						"name": "bar"
// 					}
// 				]
// 			}
// 		}`),
// 					},
// 				},
// 			},
// 			want: admission.Allowed(""),
// 		},
// 		{
// 			name: "allowed_when_device_exists_condition_missing",
// 			req: admission.Request{
// 				AdmissionRequest: admissionv1.AdmissionRequest{
// 					Object: runtime.RawExtension{
// 						Raw: []byte(`{
// 			"apiVersion": "v1",
// 			"kind": "Pod",
// 			"metadata": {
// 				"name": "foo",
// 				"namespace": "default",
// 				"annotations": {
// 					"app.kubeserial.com/inject-device": "test-device"
// 				}
// 			},
// 			"spec": {
// 				"containers": [
// 					{
// 						"image": "bar:v2",
// 						"name": "bar"
// 					}
// 				]
// 			}
// 		}`),
// 					},
// 				},
// 			},
// 			want: admission.Allowed(""),
// 		},
// 		{
// 			name: "allowed_when_device_exists",
// 			req: admission.Request{
// 				AdmissionRequest: admissionv1.AdmissionRequest{
// 					Object: runtime.RawExtension{
// 						Raw: []byte(`{
// 			"apiVersion": "v1",
// 			"kind": "Pod",
// 			"metadata": {
// 				"name": "foo",
// 				"namespace": "default",
// 				"annotations": {
// 					"app.kubeserial.com/inject-device": "test-device"
// 				}
// 			},
// 			"spec": {
// 				"containers": [
// 					{
// 						"image": "busybox:latest",
// 						"name": "bar"
// 					}
// 				]
// 			},
// 			"status": {
// 				"conditions": [
// 					{
// 						"type": "Free",
// 						"status": "True",
// 						"reason": "a",
// 						"message": "b"
// 					}
// 				]
// 			}
// 		}`),
// 					},
// 				},
// 			},
// 			want: admission.Allowed(""),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			objects := []runtime.Object{
// 				&kubeserialv1alpha1.SerialDevice{
// 					ObjectMeta: v1.ObjectMeta{
// 						Name: "test-device",
// 					},
// 				},
// 			}
// 			deviceInjector := &SerialDeviceInjector{
// 				Name:            "DeviceInjector",
// 				Clientset:       fake.NewSimpleClientset(objects...),
// 				ConfigExtractor: images.NewOCIConfigExtractor(),
// 			}
// 			decoder, err := admission.NewDecoder(scheme)
// 			assert.NoError(t, err)
// 			err = deviceInjector.InjectDecoder(decoder)
// 			assert.NoError(t, err)

// 			if got := deviceInjector.Handle(context.TODO(), tt.req); !reflect.DeepEqual(got.Allowed, tt.want.Allowed) {
// 				t.Errorf("SerialDeviceInjector.Handle() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestSerialDeviceInjector_InjectDecoder(t *testing.T) {
// 	type fields struct {
// 		Name            string
// 		Clientset       versioned.Interface
// 		ConfigExtractor *images.OCIConfigExtractor
// 		decoder         *admission.Decoder
// 	}
// 	type args struct {
// 		d *admission.Decoder
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			si := &SerialDeviceInjector{
// 				Name:            tt.fields.Name,
// 				Clientset:       tt.fields.Clientset,
// 				ConfigExtractor: tt.fields.ConfigExtractor,
// 				decoder:         tt.fields.decoder,
// 			}
// 			if err := si.InjectDecoder(tt.args.d); (err != nil) != tt.wantErr {
// 				t.Errorf("SerialDeviceInjector.InjectDecoder() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
