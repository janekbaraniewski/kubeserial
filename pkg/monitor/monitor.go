package monitor

import (
	"context"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	v1alpha1client "github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/typed/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

var log = logf.Log.WithName("DeviceMonitor")

type Monitor struct {
	cmClient      v1.ConfigMapInterface
	devicesClient v1alpha1client.SerialDeviceInterface
	namespace     string
	fs            utils.FileSystem
}

func NewMonitor(
	clientSet client.Interface,
	clientsetKubeserial versioned.Interface,
	namespace string,
	fs utils.FileSystem,
) *Monitor {
	return &Monitor{
		cmClient:      clientSet.CoreV1().ConfigMaps(namespace),
		devicesClient: clientsetKubeserial.AppV1alpha1().SerialDevices(),
		namespace:     namespace,
		fs:            fs,
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
	logger := log.WithName("updateCRDBasedDevice")
	devices, err := m.devicesClient.List(ctx, metav1.ListOptions{})
	readyDevices := []v1alpha1.SerialDevice{}
	for _, device := range devices.Items {
		readyCondition := device.GetCondition(v1alpha1.SerialDeviceReady)
		if readyCondition != nil && readyCondition.Status == metav1.ConditionTrue {
			readyDevices = append(readyDevices, device)
		}
	}
	if err != nil {
		log.Error(err, "Failed listing SerialDevice CRs")
	}
	for _, device := range readyDevices {
		logger.V(2).Info("Got device!", "device", device)
		logger = logger.WithValues("Device", device.Name)
		deviceCondition := device.GetCondition(v1alpha1.SerialDeviceAvailable)
		if deviceCondition == nil {
			log.Error(err, "Can't find device condition")
			continue
		}
		if deviceCondition.Status != metav1.ConditionTrue {
			if m.isDeviceAvailable(device.Name) {
				log.Info("Device available, updating state.")
				device.SetCondition(v1alpha1.SerialDeviceCondition{
					Type:   v1alpha1.SerialDeviceAvailable,
					Status: metav1.ConditionTrue,
					Reason: "DeviceAvailable",
				})
				device.SetCondition(v1alpha1.SerialDeviceCondition{
					Type:   v1alpha1.SerialDeviceFree,
					Status: metav1.ConditionTrue,
					Reason: "DeviceFree",
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
			device.SetCondition(v1alpha1.SerialDeviceCondition{
				Type:   v1alpha1.SerialDeviceAvailable,
				Status: metav1.ConditionFalse,
				Reason: "DeviceUnavailable",
			})
			device.SetCondition(v1alpha1.SerialDeviceCondition{
				Type:   v1alpha1.SerialDeviceFree,
				Status: metav1.ConditionUnknown,
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
	if _, err := m.fs.Stat("/dev/" + name); os.IsNotExist(err) {
		logger.V(2).Info("Device not available")
		return false
	}
	logger.V(2).Info("Device available")
	return true
}
