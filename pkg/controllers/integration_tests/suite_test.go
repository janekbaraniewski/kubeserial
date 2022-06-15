// // /*
// // Copyright 2022.

// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at

// //     http://www.apache.org/licenses/LICENSE-2.0

// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.
// // */

package integration_test

// import (
// 	"context"
// 	"math/rand"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/scheme"
// 	"k8s.io/client-go/rest"
// 	ctrl "sigs.k8s.io/controller-runtime"
// 	"sigs.k8s.io/controller-runtime/pkg/client"
// 	"sigs.k8s.io/controller-runtime/pkg/envtest"
// 	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
// 	logf "sigs.k8s.io/controller-runtime/pkg/log"
// 	"sigs.k8s.io/controller-runtime/pkg/log/zap"

// 	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
// 	"github.com/janekbaraniewski/kubeserial/pkg/controllers"
// 	//+kubebuilder:scaffold:imports
// )

// // These tests use Ginkgo (BDD-style Go testing framework). Refer to
// // http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

// var cfg *rest.Config
// var k8sClient client.Client
// var testEnv *envtest.Environment

// func TestAPIs(t *testing.T) {
// 	RegisterFailHandler(Fail)

// 	RunSpecsWithDefaultAndCustomReporters(t,
// 		"Controller Suite",
// 		[]Reporter{printer.NewlineReporter{}})
// }

// var _ = BeforeSuite(func(done Done) {
// 	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

// 	By("bootstrapping test environment")
// 	testEnv = &envtest.Environment{
// 		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "build", "_output", "kubeserial-crds", "templates")},
// 		ErrorIfCRDPathMissing: true,
// 	}

// 	cfg, err := testEnv.Start()
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(cfg).NotTo(BeNil())

// 	err = kubeserialv1alpha1.AddToScheme(scheme.Scheme)
// 	Expect(err).NotTo(HaveOccurred())

// 	//+kubebuilder:scaffold:scheme

// 	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
// 	Expect(err).NotTo(HaveOccurred())
// 	Expect(k8sClient).NotTo(BeNil())

// 	close(done)
// }, 60)

// var _ = AfterSuite(func() {
// 	By("tearing down the test environment")
// 	err := testEnv.Stop()
// 	Expect(err).NotTo(HaveOccurred())
// })

// func SetupTest(ctx context.Context) *corev1.Namespace {
// 	innerContext, cancel := context.WithCancel(ctx)
// 	ns := &corev1.Namespace{}

// 	BeforeEach(func() {
// 		*ns = corev1.Namespace{
// 			ObjectMeta: metav1.ObjectMeta{Name: "testns-" + randStringRunes(5)},
// 		}

// 		err := k8sClient.Create(ctx, ns)
// 		Expect(err).NotTo(HaveOccurred(), "failed to create test namespace")

// 		mgr, err := ctrl.NewManager(cfg, ctrl.Options{})
// 		Expect(err).NotTo(HaveOccurred(), "failed to create manager")

// 		controller := &controllers.KubeSerialReconciler{
// 			Client: mgr.GetClient(),
// 			Scheme: mgr.GetScheme(),
// 		}
// 		err = controller.SetupWithManager(mgr)
// 		Expect(err).NotTo(HaveOccurred(), "failed to setup controller")

// 		go func() {
// 			err := mgr.Start(innerContext)
// 			Expect(err).NotTo(HaveOccurred(), "failed to start manager")
// 		}()
// 	})

// 	AfterEach(func() {
// 		cancel()
// 		err := k8sClient.Delete(ctx, ns)
// 		Expect(err).NotTo(HaveOccurred(), "failed to delete test namespace")
// 	})

// 	return ns
// }

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

// func randStringRunes(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }
