//go:build e2e

/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"
	"fmt"
	"os"
	"strconv"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
)

// config holds the knobs that select where and how the suite runs. All values
// come from the environment so the same specs run locally and in CI. See the
// table in docs/e2e-testing.md.
type config struct {
	SkipClusterSetup bool
	SkipDeviceSim    bool
	KindCluster      string
	KubeContext      string
	Namespace        string
	ImageTag         string
}

var (
	cfg    config
	scheme = runtime.NewScheme()
	// k8sClient is a controller-runtime client wired with the kubeserial scheme.
	k8sClient client.Client
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func loadConfig() config {
	cfg = config{
		SkipClusterSetup: envBool("E2E_SKIP_CLUSTER_SETUP", false),
		// Device simulation (Option B: privileged socat-PTY pod on hostPath /dev)
		// is verified working on kind; enabled by default. Set true to skip the
		// device specs (E4/E5) on substrates where hostPath /dev is unavailable.
		SkipDeviceSim: envBool("E2E_SKIP_DEVICE_SIM", false),
		KindCluster:   envOr("E2E_KIND_CLUSTER", "kubeserial-e2e"),
		KubeContext:   envOr("E2E_KUBECONTEXT", "kind-kubeserial-e2e"),
		Namespace:     envOr("E2E_NAMESPACE", "kubeserial"),
		ImageTag:      envOr("E2E_IMAGE_TAG", "local"),
	}
	return cfg
}

// restConfigForContext builds a *rest.Config. In CI/kind the kubeconfig current
// context is set by the workflow; locally we honor E2E_KUBECONTEXT.
func restConfigForContext(kubeContext string) (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
}

func initClient(c config) {
	restCfg, err := restConfigForContext(c.KubeContext)
	Expect(err).NotTo(HaveOccurred(), "failed to build rest.Config for context %q", c.KubeContext)

	cl, err := client.New(restCfg, client.Options{Scheme: scheme})
	Expect(err).NotTo(HaveOccurred(), "failed to construct controller-runtime client")
	k8sClient = cl
}

// ensureChartInstalled is a light precondition check used when the suite did not
// itself install the chart. It verifies the SerialDevice CRD is registered.
func ensureChartInstalled(_ config) {
	Eventually(func() error {
		list := &kubeserialv1alpha1.SerialDeviceList{}
		return k8sClient.List(context.Background(), list)
	}).Should(Succeed(), "SerialDevice CRD should be installed by the kubeserial chart")
}

// ---------------------------------------------------------------------------
// SerialDevice helpers
// ---------------------------------------------------------------------------

// newSerialDevice returns a minimal SerialDevice CR. The monitor stats
// /dev/<metadata.name>, so the name doubles as the device-node name.
func newSerialDevice(name string) *kubeserialv1alpha1.SerialDevice {
	return &kubeserialv1alpha1.SerialDevice{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec: kubeserialv1alpha1.SerialDeviceSpec{
			Name:      name,
			IDVendor:  "0403",
			IDProduct: "6001",
		},
	}
}

// conditionStatus returns the status of the named condition, or "" if absent.
func conditionStatus(
	dev *kubeserialv1alpha1.SerialDevice,
	t kubeserialv1alpha1.SerialDeviceConditionType,
) metav1.ConditionStatus {
	if c := dev.GetCondition(t); c != nil {
		return c.Status
	}
	return ""
}

// getDevice fetches a fresh copy of a SerialDevice by name.
func getDevice(ctx context.Context, name string) (*kubeserialv1alpha1.SerialDevice, error) {
	dev := &kubeserialv1alpha1.SerialDevice{}
	err := k8sClient.Get(ctx, client.ObjectKey{Name: name}, dev)
	return dev, err
}

// setDeviceCondition patches a condition onto a SerialDevice's status. Used by
// webhook specs (E6/E7) to put a device into a known Free state without real
// hardware. It retries on conflict because the SerialDevice controller
// reconciles the same object concurrently (e.g. setting Ready), which otherwise
// races our read-modify-write and yields a 409.
func setDeviceCondition(
	ctx context.Context,
	name string,
	cond kubeserialv1alpha1.SerialDeviceCondition,
) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		dev, err := getDevice(ctx, name)
		if err != nil {
			return err
		}
		dev.SetCondition(cond)
		return k8sClient.Status().Update(ctx, dev)
	})
}

// ---------------------------------------------------------------------------
// Device-simulation harness (Option B). See docs/e2e-testing.md section 2.
// ---------------------------------------------------------------------------

// simulatorImage is the container image used by the device-simulator pod. It
// must contain `socat` and a shell. alpine/socat satisfies both. hack/e2e.sh
// pre-loads this image into the kind node so the pod starts without registry
// access.
const simulatorImage = "alpine/socat:latest"

// deviceSimulatorPod returns a privileged pod that creates a PTY-backed char
// device symlinked at /dev/<deviceName> on the node's host /dev (which the
// monitor DaemonSet mounts via hostPath), then blocks. Deleting the pod closes
// the PTY, so the symlink target disappears and the monitor's os.Stat (which
// follows the symlink) reports the device gone, simulating unplug.
//
// VERIFIED on kind v0.18 (kindest/node v1.27.1): the device created here is
// visible inside the monitor pod via the shared hostPath /dev with no extra
// mount-propagation handling, and the monitor flips Available/Free accordingly.
func deviceSimulatorPod(namespace, deviceName string) *corev1.Pod {
	privileged := true
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "device-sim-" + deviceName,
			Namespace: namespace,
			Labels:    map[string]string{"app.kubeserial.com/e2e": "device-sim"},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{{
				Name:            "socat",
				Image:           simulatorImage,
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: &corev1.SecurityContext{Privileged: &privileged},
				Command:         []string{"/bin/sh", "-c"},
				Args: []string{fmt.Sprintf(
					// Create a PTY char device symlinked at /dev/<name>, then block.
					"socat -d -d PTY,raw,echo=0,link=/dev/%s PTY,raw,echo=0 & "+
						"while true; do sleep 3600; done",
					deviceName,
				)},
				VolumeMounts: []corev1.VolumeMount{{Name: "host-dev", MountPath: "/dev"}},
			}},
			Volumes: []corev1.Volume{{
				Name: "host-dev",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/dev"},
				},
			}},
		},
	}
}

// simulateDeviceAttach schedules the simulator pod so /dev/<name> appears on the
// node. simulateDeviceDetach deletes it so the device disappears.
func simulateDeviceAttach(ctx context.Context, namespace, deviceName string) {
	pod := deviceSimulatorPod(namespace, deviceName)
	Expect(k8sClient.Create(ctx, pod)).To(Succeed(), "create device simulator pod")
	Eventually(func() (corev1.PodPhase, error) {
		got := &corev1.Pod{}
		err := k8sClient.Get(ctx, client.ObjectKeyFromObject(pod), got)
		return got.Status.Phase, err
	}).Should(Equal(corev1.PodRunning), "device simulator pod should be Running")
}

func simulateDeviceDetach(ctx context.Context, namespace, deviceName string) {
	pod := deviceSimulatorPod(namespace, deviceName)
	_ = k8sClient.Delete(ctx, pod)
	// Block until the pod is fully gone so the PTY is closed before callers
	// assert the device disappeared.
	Eventually(func() bool {
		got := &corev1.Pod{}
		err := k8sClient.Get(ctx, client.ObjectKeyFromObject(pod), got)
		return apierrors.IsNotFound(err)
	}).Should(BeTrue(), "device simulator pod should be deleted")
}
