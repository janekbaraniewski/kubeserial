package main

import (
	"flag"
	"os"

	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/images"
	"github.com/janekbaraniewski/kubeserial/pkg/webhooks"
)

var setupLog = ctrl.Log.WithName("setup")

type HookParamters struct {
	certDir   string
	port      int
	namespace string
}

func main() {
	var params HookParamters

	flag.IntVar(&params.port, "port", 8443, "Wehbook port")
	flag.StringVar(&params.certDir, "certDir", "/certs/", "Wehbook certificate folder")
	flag.StringVar(&params.namespace, "namespace", "kubeserial", "Namespace in which controller is deployed") // TODO: remove this and make it work globaly
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	entryLog := setupLog.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		MetricsBindAddress:     ":8080",
		HealthProbeBindAddress: ":8081",
	})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	if err != nil {
		entryLog.Error(err, "Can't load config")
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Port = params.port
	hookServer.CertDir = params.certDir

	config, err := rest.InClusterConfig()

	if err != nil {
		entryLog.Error(err, "Failed to get InClusterConfig")
		panic(err.Error())
	}

	clientset, err := versioned.NewForConfig(config)

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register(
		"/mutate-inject-device",
		&webhook.Admission{
			Handler: &webhooks.DeviceInjector{
				Name:            "DeviceInjector",
				Clientset:       clientset,
				ConfigExtractor: images.NewOCIConfigExtractor(),
				Namespace:       params.namespace,
			},
		},
	)

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
