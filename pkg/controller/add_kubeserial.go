package controller

import (
	"github.com/janekbaraniewski/kubeserial/pkg/controller/kubeserial"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, kubeserial.Add)
}
