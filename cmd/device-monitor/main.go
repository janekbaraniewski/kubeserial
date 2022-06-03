package main

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	log := ctrl.Log.WithName("monitor")

	log.Info("Start setup")

	_, err := metrics.NewListener(":8080")

	if err != nil {
		log.Info("Failed setting up metrics listener")
	}

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
	deviceMonitor := monitor.NewMonitor(
		clientset,
		clientsetKubeserial,
		os.Getenv("OPERATOR_NAMESPACE"),
		utils.NewOSFS(),
	)
	log.Info("Starting monitor update loop")
	go deviceMonitor.RunUpdateLoop(ctx)
	<-ctx.Done()
	log.Info("Exiting")
}
