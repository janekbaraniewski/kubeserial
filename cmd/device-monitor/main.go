package main

import (
	"os"

	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	log := ctrl.Log.WithName("monitor")

	log.Info("Start setup")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error(err, "Failed to get InClusterConfig")
		panic(err.Error())
	}

	clientsetKubeserial, err := versioned.NewForConfig(config)
	if err != nil {
		log.Error(err, "Can't create kubeserial clientset")
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error(err, "Failed to get clientset")
		panic(err.Error())
	}
	log.Info("Clientset initialised")
	ctx := ctrl.SetupSignalHandler()
	log.Info("Starting update loop")
	go monitor.RunUpdateLoop(
		ctx,
		clientset,
		os.Getenv("OPERATOR_NAMESPACE"),
		clientsetKubeserial,
	)
	<-ctx.Done()
	log.Info("Exiting")
}
