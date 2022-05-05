package integration_tests

// import (
// 	"context"
// 	"reflect"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	batchv1 "k8s.io/api/batch/v1"
// 	batchv1beta1 "k8s.io/api/batch/v1beta1"
// 	v1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/types"
// )

// var _ = Describe("CronJob controller", func() {

// 	// Define utility constants for object names and testing timeouts/durations and intervals.
// 	const (
// 		CronjobName      = "test-cronjob"
// 		CronjobNamespace = "default"
// 		JobName          = "test-job"

// 		timeout  = time.Second * 10
// 		duration = time.Second * 10
// 		interval = time.Millisecond * 250
// 	)

// 	Context("When updating CronJob Status", func() {
// 		It("Should increase CronJob Status.Active count when new Jobs are created", func() {
// 			By("By creating a new CronJob")
// 			ctx := context.Background()
// 			cronJob := &cronjobv1.CronJob{
// 				TypeMeta: metav1.TypeMeta{
// 					APIVersion: "batch.tutorial.kubebuilder.io/v1",
// 					Kind:       "CronJob",
// 				},
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      CronjobName,
// 					Namespace: CronjobNamespace,
// 				},
// 				Spec: cronjobv1.CronJobSpec{
// 					Schedule: "1 * * * *",
// 					JobTemplate: batchv1beta1.JobTemplateSpec{
// 						Spec: batchv1.JobSpec{
// 							// For simplicity, we only fill out the required fields.
// 							Template: v1.PodTemplateSpec{
// 								Spec: v1.PodSpec{
// 									// For simplicity, we only fill out the required fields.
// 									Containers: []v1.Container{
// 										{
// 											Name:  "test-container",
// 											Image: "test-image",
// 										},
// 									},
// 									RestartPolicy: v1.RestartPolicyOnFailure,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			}
// 			Expect(k8sClient.Create(ctx, cronJob)).Should(Succeed())

// 			cronjobLookupKey := types.NamespacedName{Name: CronjobName, Namespace: CronjobNamespace}
// 			createdCronjob := &cronjobv1.CronJob{}

// 			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
// 			Eventually(func() bool {
// 				err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
// 				if err != nil {
// 					return false
// 				}
// 				return true
// 			}, timeout, interval).Should(BeTrue())
// 			// Let's make sure our Schedule string value was properly converted/handled.
// 			Expect(createdCronjob.Spec.Schedule).Should(Equal("1 * * * *"))

// 			By("By checking the CronJob has zero active Jobs")
// 			Consistently(func() (int, error) {
// 				err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
// 				if err != nil {
// 					return -1, err
// 				}
// 				return len(createdCronjob.Status.Active), nil
// 			}, duration, interval).Should(Equal(0))

// 			By("By creating a new Job")
// 			testJob := &batchv1.Job{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      JobName,
// 					Namespace: CronjobNamespace,
// 				},
// 				Spec: batchv1.JobSpec{
// 					Template: v1.PodTemplateSpec{
// 						Spec: v1.PodSpec{
// 							// For simplicity, we only fill out the required fields.
// 							Containers: []v1.Container{
// 								{
// 									Name:  "test-container",
// 									Image: "test-image",
// 								},
// 							},
// 							RestartPolicy: v1.RestartPolicyOnFailure,
// 						},
// 					},
// 				},
// 				Status: batchv1.JobStatus{
// 					Active: 2,
// 				},
// 			}

// 			// Note that your CronJob’s GroupVersionKind is required to set up this owner reference.
// 			kind := reflect.TypeOf(cronjobv1.CronJob{}).Name()
// 			gvk := cronjobv1.GroupVersion.WithKind(kind)

// 			controllerRef := metav1.NewControllerRef(createdCronjob, gvk)
// 			testJob.SetOwnerReferences([]metav1.OwnerReference{*controllerRef})
// 			Expect(k8sClient.Create(ctx, testJob)).Should(Succeed())

// 			By("By checking that the CronJob has one active Job")
// 			Eventually(func() ([]string, error) {
// 				err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
// 				if err != nil {
// 					return nil, err
// 				}

// 				names := []string{}
// 				for _, job := range createdCronjob.Status.Active {
// 					names = append(names, job.Name)
// 				}
// 				return names, nil
// 			}, timeout, interval).Should(ConsistOf(JobName), "should list our active job %s in the active jobs list in status", JobName)
// 		})
// 	})

// })
