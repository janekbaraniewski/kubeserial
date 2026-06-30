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
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
)

// Webhook specs (E6/E7) need NO real hardware: the webhook decides based on the
// SerialDevice's Free condition, which we set directly on the status subresource.
// This is the highest-value path that is fully testable in hosted CI.
//
// These are scaffolded but gated behind E2E_SKIP_DEVICE_SIM only insofar as they
// require the webhook + its MutatingWebhookConfiguration to be installed and the
// TLS/caBundle wiring to work under kind, which is UNVERIFIED here. They run
// whenever the webhook is present; if pod creation never gets mutated because
// the webhook is not wired, E6 fails loudly (which is the signal to fix wiring).
const (
	injectAnnotation = "app.kubeserial.com/inject-device"
	socatMarker      = "socat"
	gatewaySuffix    = "-gateway:3333"
)

var _ = Describe("device-injection webhook", func() {
	ctx := context.Background()

	makePod := func(name, device string) *corev1.Pod {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   cfg.Namespace,
				Annotations: map[string]string{injectAnnotation: device},
			},
			Spec: corev1.PodSpec{
				// A command is required so the webhook has something to wrap; an
				// empty command would force the webhook down the image-config
				// extraction path which needs registry access.
				Containers: []corev1.Container{{
					Name:    "app",
					Image:   "busybox:latest",
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", "sleep 3600"},
				}},
			},
		}
	}

	// E6: when the requested device is Free=True, the webhook rewrites
	// containers[0] to wrap the original command with a socat bridge to the
	// device gateway.
	It("injects a socat bridge when the device is Free [E6]", func() {
		const deviceName = "e2e-webhook-free"
		dev := newSerialDevice(deviceName)
		Expect(k8sClient.Create(ctx, dev)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, dev) })

		Expect(setDeviceCondition(ctx, deviceName, kubeserialv1alpha1.SerialDeviceCondition{
			Type:   kubeserialv1alpha1.SerialDeviceFree,
			Status: metav1.ConditionTrue,
			Reason: "E2ESetup",
		})).To(Succeed())

		pod := makePod("e2e-inject-yes", deviceName)
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, pod) })

		got := &corev1.Pod{}
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(pod), got)).To(Succeed())

		joined := strings.Join(append(got.Spec.Containers[0].Command, got.Spec.Containers[0].Args...), " ")
		Expect(joined).To(ContainSubstring(socatMarker),
			"expected webhook to inject a socat bridge into the container command")
		Expect(joined).To(ContainSubstring(deviceName+gatewaySuffix),
			"expected the socat bridge to target the device gateway")
	})

	// E7: when the device is NOT Free, the webhook must leave the pod unchanged.
	It("does not inject when the device is not Free [E7]", func() {
		const deviceName = "e2e-webhook-busy"
		dev := newSerialDevice(deviceName)
		Expect(k8sClient.Create(ctx, dev)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, dev) })

		// Explicitly mark Free=False.
		Expect(setDeviceCondition(ctx, deviceName, kubeserialv1alpha1.SerialDeviceCondition{
			Type:   kubeserialv1alpha1.SerialDeviceFree,
			Status: metav1.ConditionFalse,
			Reason: "E2ESetup",
		})).To(Succeed())

		pod := makePod("e2e-inject-no", deviceName)
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, pod) })

		got := &corev1.Pod{}
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(pod), got)).To(Succeed())

		joined := strings.Join(append(got.Spec.Containers[0].Command, got.Spec.Containers[0].Args...), " ")
		Expect(joined).NotTo(ContainSubstring(socatMarker),
			"webhook must not inject a socat bridge when the device is not Free")
	})

	// E7b: a pod with no inject-device annotation must be left untouched, even
	// though the webhook intercepts every pod CREATE.
	It("does not inject when the pod has no inject annotation [E7]", func() {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "e2e-inject-none",
				Namespace: cfg.Namespace,
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
		Expect(k8sClient.Create(ctx, pod)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, pod) })

		got := &corev1.Pod{}
		Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(pod), got)).To(Succeed())

		joined := strings.Join(append(got.Spec.Containers[0].Command, got.Spec.Containers[0].Args...), " ")
		Expect(joined).NotTo(ContainSubstring(socatMarker),
			"webhook must not touch a pod without the inject-device annotation")
		Expect(joined).To(ContainSubstring("sleep 3600"),
			"original command must be preserved")
	})
})
