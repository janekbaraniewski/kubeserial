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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
)

// E4/E5 use the device-simulation harness (a privileged pod that creates
// /dev/<name> on the node); set E2E_SKIP_DEVICE_SIM=true to skip them on
// substrates without hostPath /dev.
var _ = Describe("device-monitor presence detection", Ordered, func() {
	ctx := context.Background()
	const deviceName = "e2e-sim-device"

	BeforeEach(func() {
		if cfg.SkipDeviceSim {
			Skip("E2E_SKIP_DEVICE_SIM is true")
		}
	})

	var dev *kubeserialv1alpha1.SerialDevice

	BeforeAll(func() {
		if cfg.SkipDeviceSim {
			return
		}
		dev = newSerialDevice(deviceName)
		Expect(k8sClient.Create(ctx, dev)).To(Succeed())
		// The monitor only evaluates presence for Ready devices.
		Expect(setDeviceCondition(ctx, deviceName, kubeserialv1alpha1.SerialDeviceCondition{
			Type:   kubeserialv1alpha1.SerialDeviceReady,
			Status: metav1.ConditionTrue,
			Reason: "E2ESetup",
		})).To(Succeed())
	})

	AfterAll(func() {
		if cfg.SkipDeviceSim {
			return
		}
		simulateDeviceDetach(ctx, cfg.Namespace, deviceName)
		if dev != nil {
			_ = k8sClient.Delete(ctx, dev)
		}
	})

	It("flips Available and Free to True when the device appears [E4]", func() {
		simulateDeviceAttach(ctx, cfg.Namespace, deviceName)

		Eventually(func(g Gomega) {
			got, err := getDevice(ctx, deviceName)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(conditionStatus(got, kubeserialv1alpha1.SerialDeviceAvailable)).
				To(Equal(metav1.ConditionTrue))
			g.Expect(conditionStatus(got, kubeserialv1alpha1.SerialDeviceFree)).
				To(Equal(metav1.ConditionTrue))
			g.Expect(got.Status.NodeName).NotTo(BeEmpty())
		}).Should(Succeed())
	})

	It("flips Available back to False when the device disappears [E5]", func() {
		simulateDeviceDetach(ctx, cfg.Namespace, deviceName)

		Eventually(func(g Gomega) {
			got, err := getDevice(ctx, deviceName)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(conditionStatus(got, kubeserialv1alpha1.SerialDeviceAvailable)).
				To(Equal(metav1.ConditionFalse))
			g.Expect(got.Status.NodeName).To(BeEmpty())
		}).Should(Succeed())
	})
})
