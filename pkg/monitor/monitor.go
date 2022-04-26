package monitor

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	client "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("ApiClient")

func RunUpdateLoop(ctx context.Context, clientset client.Interface, namespace string, clientsetKubeserial versioned.Interface) {
	for {
		select {
		case <-time.After(1 * time.Second):
			UpdateDeviceState(ctx, clientset, clientsetKubeserial, namespace)
		case <-ctx.Done():
			return
		}
	}
}

func UpdateDeviceState(ctx context.Context, clientset client.Interface, clientsetKubeserial versioned.Interface, namespace string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		updateCMBasedDevice(ctx, clientset.CoreV1().ConfigMaps(namespace))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		UpdateCRDBasedDevice(ctx, clientsetKubeserial, namespace)
	}()

	wg.Wait()
}

func updateCMBasedDevice(ctx context.Context, client v1.ConfigMapInterface) {
	confs, err := client.List(ctx, metav1.ListOptions{
		LabelSelector: "type=DeviceState", // TODO: make configurable
	})
	if err != nil {
		panic(err.Error())
	}

	for _, conf := range confs.Items {
		if conf.Data["node"] == os.Getenv("NODE_NAME") {
			if !isDeviceAvailable(conf.Labels["device"]) {
				log.Info("Device unavailable, cleaning state.")
				if err := clearState(ctx, &conf, client); err != nil {
					log.Error(err, "Update failed to clear state!")
				}
			}
		} else if conf.Data["available"] == "false" {
			if isDeviceAvailable(conf.Labels["device"]) {
				log.Info("Device available, updating state.")
				if err := setActiveState(ctx, &conf, client); err != nil {
					log.Error(err, "Update failed to make device available!")
				}
			}
		}
	}

}

func UpdateCRDBasedDevice(ctx context.Context, clientset versioned.Interface, namespace string) {
	client := clientset.AppV1alpha1().Devices(namespace)

	devices, err := client.List(ctx, metav1.ListOptions{
		FieldSelector: fields.SelectorFromSet(fields.Set{
			"status.conditions[].type":   string(v1alpha1.DeviceReady),
			"status.conditions[].status": string(corev1.ConditionTrue),
		}).String(),
	})
	if err != nil {
		log.Error(err, "Failed listing Device CRs")
	}
	for _, device := range devices.Items {
		log.V(2).Info("Got device!", "device", device)
		deviceCondition := utils.GetCondition(device.Status.Conditions, v1alpha1.DeviceAvailable)
		if deviceCondition == nil {
			log.Error(err, "Can't find device condition")
			continue
		}
		if deviceCondition.Status == metav1.ConditionFalse {
			if isDeviceAvailable(device.Name) {
				log.Info("Device available, updating state.")
				//TODO: update condition
			}
		}
	}
}

func isDeviceAvailable(name string) bool {
	if _, err := os.Stat("/dev/tty" + name); os.IsNotExist(err) {
		return false
	}
	return true
}

func clearState(ctx context.Context, c *corev1.ConfigMap, client v1.ConfigMapInterface) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "false"
		c.Data["node"] = ""
		_, updateErr := client.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}

func setActiveState(ctx context.Context, c *corev1.ConfigMap, client v1.ConfigMapInterface) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "true"
		c.Data["node"] = os.Getenv("NODE_NAME")
		_, updateErr := client.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}
