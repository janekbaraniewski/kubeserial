package monitor

import (
	"context"
	"os"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	v1alpha1client "github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/typed/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/metrics"
)

var log = logf.Log.WithName("DeviceMonitor")

type Monitor struct {
	cmClient      v1.ConfigMapInterface
	devicesClient v1alpha1client.DeviceInterface
	namespace     string
	statFile      func(filename string) (os.FileInfo, error)
}

func NewMonitor(clientSet client.Interface, clientsetKubeserial versioned.Interface, namespace string, statFunc func(filename string) (os.FileInfo, error)) *Monitor {
	return &Monitor{
		cmClient:      clientSet.CoreV1().ConfigMaps(namespace),
		devicesClient: clientsetKubeserial.AppV1alpha1().Devices(namespace),
		namespace:     namespace,
		statFile:      statFunc,
	}
}

func (m *Monitor) RunUpdateLoop(ctx context.Context) {
	log.Info("Starting update loop")
	for {
		select {
		case <-time.After(1 * time.Second):
			m.UpdateDeviceState(ctx)
		case <-ctx.Done():
			log.Info("Stopping update loop")
			return
		}
	}
}

func (m *Monitor) UpdateDeviceState(ctx context.Context) {
	logger := log.WithName("UpdateDeviceState")
	logger.V(2).Info("Updating device state")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.updateCMBasedDevice(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.updateCRDBasedDevice(ctx)
	}()

	wg.Wait()
}

func (m *Monitor) updateCMBasedDevice(ctx context.Context) {
	confs, err := m.cmClient.List(ctx, metav1.ListOptions{
		LabelSelector: "type=DeviceState", // TODO: make configurable
	})
	if err != nil {
		panic(err.Error())
	}

	for _, conf := range confs.Items {
		if conf.Data["node"] == os.Getenv("NODE_NAME") {
			if !m.isDeviceAvailable(conf.Labels["device"]) {
				log.Info("Device unavailable, cleaning state.")
				metrics.AvailableDevices.Dec()
				if err := m.clearState(ctx, &conf); err != nil {
					log.Error(err, "Update failed to clear state!")
				}
			}
		} else if conf.Data["available"] == "false" {
			if m.isDeviceAvailable(conf.Labels["device"]) {
				log.Info("Device available, updating state.")
				metrics.AvailableDevices.Inc()
				if err := m.setActiveState(ctx, &conf); err != nil {
					log.Error(err, "Update failed to make device available!")
				}
			}
		}
	}

}

func (m *Monitor) updateCRDBasedDevice(ctx context.Context) {
	logger := log.WithName("updateCRDBasedDevice")
	devices, err := m.devicesClient.List(ctx, metav1.ListOptions{})
	readyDevices := []v1alpha1.Device{}
	for _, device := range devices.Items {
		readyCondition := device.GetCondition(v1alpha1.DeviceReady)
		if readyCondition != nil && readyCondition.Status == metav1.ConditionTrue {
			readyDevices = append(readyDevices, device)
		}
	}
	if err != nil {
		log.Error(err, "Failed listing Device CRs")
	}
	for _, device := range readyDevices {
		logger.V(2).Info("Got device!", "device", device)
		logger = logger.WithValues("Device", device.Name)
		deviceCondition := device.GetCondition(v1alpha1.DeviceAvailable)
		if deviceCondition == nil {
			log.Error(err, "Can't find device condition")
			continue
		}
		if deviceCondition.Status != metav1.ConditionTrue {
			if m.isDeviceAvailable(device.Name) {
				log.Info("Device available, updating state.")
				device.SetCondition(v1alpha1.DeviceCondition{
					Type:   v1alpha1.DeviceAvailable,
					Status: metav1.ConditionTrue,
					Reason: "DeviceAvailable",
				})
				device.Status.NodeName = os.Getenv("NODE_NAME")
				logger.WithValues("Node", device.Status.NodeName).Info("Setting device state to available")
				_, err := m.devicesClient.UpdateStatus(ctx, &device, metav1.UpdateOptions{})
				if err != nil {
					log.Error(err, "Failed device status update")
				}
			}
		} else if device.Status.NodeName == os.Getenv("NODE_NAME") && !m.isDeviceAvailable(device.Name) {
			log.Info("Device unavailable, updating state.")
			device.SetCondition(v1alpha1.DeviceCondition{
				Type:   v1alpha1.DeviceAvailable,
				Status: metav1.ConditionFalse,
				Reason: "DeviceUnavailable",
			})
			device.Status.NodeName = ""
			logger.Info("Setting device state to unavailable")
			_, err := m.devicesClient.UpdateStatus(ctx, &device, metav1.UpdateOptions{})
			if err != nil {
				log.Error(err, "Failed device status update")
			}
		}
	}
}

func (m *Monitor) isDeviceAvailable(name string) bool {
	logger := log.WithName("isDeviceAvailable").WithValues("Device", name)
	if _, err := m.statFile("/dev/" + name); os.IsNotExist(err) {
		logger.V(2).Info("Device not available")
		return false
	}
	logger.V(2).Info("Device available")
	return true
}

func (m *Monitor) clearState(ctx context.Context, c *corev1.ConfigMap) error {
	log.WithName("setActiveState").WithValues("Device", c.Labels["device"]).Info("[Legacy mode] Clearing device state to active")
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "false"
		c.Data["node"] = ""
		_, updateErr := m.cmClient.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}

func (m *Monitor) setActiveState(ctx context.Context, c *corev1.ConfigMap) error {
	log.WithName("setActiveState").WithValues("Device", c.Labels["device"]).Info("[Legacy mode] Setting device state to active")
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		c.Data["available"] = "true"
		c.Data["node"] = os.Getenv("NODE_NAME")
		_, updateErr := m.cmClient.Update(ctx, c, metav1.UpdateOptions{})
		return updateErr
	})
}
