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

// Package e2e runs kubeserial against a real cluster (kind in CI) with the Helm
// charts installed. See docs/e2e-testing.md.
package e2e

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
	initClient(cfg)
	if cfg.SkipClusterSetup {
		return
	}
	ensureChartInstalled(cfg)
})

// Cluster create/install/teardown is driven by `make test-e2e` and the CI
// workflow, not from inside the suite, so logs can be exported before deletion.
