package utils

import (
	"fmt"
	"time"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrConditionNotFound = fmt.Errorf("condition type not found")

func GetCondition(conditions []v1alpha1.DeviceCondition, conditionType v1alpha1.DeviceConditionType) *v1alpha1.DeviceCondition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

func SetDeviceCondition(conditions *[]v1alpha1.DeviceCondition, newCondition v1alpha1.DeviceCondition) {
	existing := GetCondition(*conditions, newCondition.Type)

	if existing == nil {
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = v1.NewTime(time.Now())
		}
		newCondition.LastHeartbeatTime = v1.NewTime(time.Now())
		*conditions = append(*conditions, newCondition)
		return
	}

	if existing.Status != newCondition.Status {
		existing.Status = newCondition.Status
		if !newCondition.LastTransitionTime.IsZero() {
			existing.LastTransitionTime = newCondition.LastTransitionTime
		} else {
			existing.LastTransitionTime = v1.NewTime(time.Now())
		}
	}

	existing.Reason = newCondition.Reason
	existing.Message = newCondition.Message
	existing.ObservedGeneration = newCondition.ObservedGeneration
	existing.LastHeartbeatTime = v1.NewTime(time.Now())
}
