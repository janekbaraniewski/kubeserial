package main

import (
	"context"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	clientCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{}
	log := zap.New(zap.UseFlagOptions(&opts))

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Failed to get InClusterConfig")
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Failed to get clientset")
		panic(err.Error())
	}
	client := clientset.CoreV1().ConfigMaps(os.Getenv("OPERATOR_NAMESPACE"))
	log.Info("Client initialised")
	ctx := ctrl.SetupSignalHandler()
	for {
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

		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func isDeviceAvailable(name string) bool {
	if _, err := os.Stat("/dev/tty" + name); os.IsNotExist(err) {
		return false
	}
	return true
}

func clearState(ctx context.Context, c *corev1.ConfigMap, client clientCorev1.ConfigMapInterface) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "false"
		c.Data["node"] = ""
		_, updateErr := client.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}

func setActiveState(ctx context.Context, c *corev1.ConfigMap, client clientCorev1.ConfigMapInterface) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "true"
		c.Data["node"] = os.Getenv("NODE_NAME")
		_, updateErr := client.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}
