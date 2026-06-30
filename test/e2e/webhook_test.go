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

const (
	injectAnnotation = "app.kubeserial.com/inject-device"
	socatMarker      = "socat"
	gatewaySuffix    = "-gateway:3333"
)

var _ = Describe("device-injection webhook", func() {
	ctx := context.Background()

	// A command is required so the webhook wraps it rather than taking the
	// image-config extraction path (which needs registry access).
	makePod := func(name, device string) *corev1.Pod {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   cfg.Namespace,
				Annotations: map[string]string{injectAnnotation: device},
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

	It("does not inject when the device is not Free [E7]", func() {
		const deviceName = "e2e-webhook-busy"
		dev := newSerialDevice(deviceName)
		Expect(k8sClient.Create(ctx, dev)).To(Succeed())
		DeferCleanup(func() { _ = k8sClient.Delete(ctx, dev) })

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
