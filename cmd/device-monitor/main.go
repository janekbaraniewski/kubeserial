package main

import (
	"flag"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/monitor"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
)

func main() {
	var (
		namespace string
		nodeName  string
	)
	flag.StringVar(&namespace, "namespace", os.Getenv("OPERATOR_NAMESPACE"), "Namespace the device monitor watches")
	flag.StringVar(&nodeName, "node-name", os.Getenv("NODE_NAME"), "Name of the node this monitor runs on")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

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
	deviceMonitor := monitor.NewMonitor(
		clientset,
		clientsetKubeserial,
		namespace,
		nodeName,
		utils.NewOSFS(),
	)
	log.Info("Starting monitor update loop")
	go deviceMonitor.RunUpdateLoop(ctx)
	<-ctx.Done()
	log.Info("Exiting")
}
