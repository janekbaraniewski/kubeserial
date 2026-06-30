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

// Package e2e contains the end-to-end test suite for kubeserial. It runs
// against a real Kubernetes cluster (kind in CI) with the kubeserial Helm
// charts installed. See docs/e2e-testing.md for the design and the
// device-simulation strategy.
//
// Unlike the unit tests (controller-runtime fake client + afero in-memory FS)
// and the envtest integration tests, this suite exercises behaviors that only
// exist on a real node with a real kubelet and a real /dev: the device-monitor
// DaemonSet, RBAC, and the mutating webhook admission round-trip.
package e2e

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Suite-wide polling defaults. The monitor polls /dev every 1s and the
// controllers requeue, so generous Eventually windows avoid flakes.
const (
	pollInterval = 2 * time.Second
	pollTimeout  = 2 * time.Minute
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyPollingInterval(pollInterval)
	SetDefaultEventuallyTimeout(pollTimeout)
	RunSpecs(t, "kubeserial e2e suite")
}

var _ = BeforeSuite(func() {
	cfg := loadConfig()

	By("connecting a controller-runtime client to the target cluster")
	initClient(cfg)

	if cfg.SkipClusterSetup {
		By("E2E_SKIP_CLUSTER_SETUP=true: assuming cluster + chart are already installed")
		return
	}

	// NOTE: cluster creation, image build/load, and `helm install` are driven by
	// the Makefile target `test-e2e` and the GitHub Actions workflow, NOT from
	// inside the suite. This keeps the Go suite focused on assertions and lets
	// the same specs run against any prepared cluster (kind, k3d, real). If you
	// prefer in-suite setup, wire installChart()/createCluster() helpers here.
	By("verifying the kubeserial chart is installed (CRDs + manager present)")
	ensureChartInstalled(cfg)
})

var _ = AfterSuite(func() {
	// Cluster teardown is the Makefile/CI's responsibility (so logs can be
	// exported on failure before deletion). Nothing to do here unless we created
	// per-spec namespaces, which are cleaned up in the specs themselves.
})
