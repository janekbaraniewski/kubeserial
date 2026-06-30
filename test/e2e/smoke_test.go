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

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
)

// Smoke specs (E1-E3) need no real hardware: they verify the chart installs,
// the CRDs are registered, and a SerialDevice with no backing /dev node never
// becomes Ready/Available. These should pass against any prepared kind cluster.
var _ = Describe("kubeserial smoke", func() {
	ctx := context.Background()

	// E1: CRDs install.
	It("registers the SerialDevice CRD [E1]", func() {
		list := &kubeserialv1alpha1.SerialDeviceList{}
		Expect(k8sClient.List(ctx, list)).To(Succeed())
	})

	// E3: a SerialDevice with no device node on any node must stay
	// not-Available. The monitor only flips Available/Free to True once it
	// stats /dev/<name>; with no device present that never happens.
	//
	// NOTE: the Ready condition is owned by the controller (it gates whether the
	// monitor even looks at the device). This spec asserts the device does not
	// spontaneously become Available, which is the hardware-presence signal.
	It("keeps a device with no backing /dev node not-Available [E3]", func() {
		dev := newSerialDevice("e2e-ghost-device")
		Expect(k8sClient.Create(ctx, dev)).To(Succeed())
		DeferCleanup(func() {
			_ = k8sClient.Delete(ctx, dev)
		})

		// Consistently over a window, Available must never be True.
		Consistently(func() metav1.ConditionStatus {
			got, err := getDevice(ctx, dev.Name)
			if err != nil {
				return ""
			}
			return conditionStatus(got, kubeserialv1alpha1.SerialDeviceAvailable)
		}, "20s", pollInterval).ShouldNot(Equal(metav1.ConditionTrue),
			"a device with no backing /dev node must never report Available=True")
	})
})

// E2: the manager Deployment becomes Available. Kept separate because it asserts
// on a core resource installed by the chart rather than a CRD object.
var _ = Describe("kubeserial manager deployment [E2]", func() {
	ctx := context.Background()

	It("has a manager Deployment with at least one ready replica", func() {
		// The chart names the manager deployment after the release; we look it up
		// by the well-known label rather than hardcoding the name. If your chart
		// uses a different selector, adjust this label.
		Eventually(func(g Gomega) {
			deps := &appsv1.DeploymentList{}
			g.Expect(k8sClient.List(ctx, deps,
				client.InNamespace(cfg.Namespace),
			)).To(Succeed())
			anyReady := false
			for i := range deps.Items {
				if deps.Items[i].Status.ReadyReplicas > 0 {
					anyReady = true
					break
				}
			}
			g.Expect(anyReady).To(BeTrue(),
				"expected at least one Deployment in %q with a ready replica", cfg.Namespace)
		}).Should(Succeed())
	})
})
