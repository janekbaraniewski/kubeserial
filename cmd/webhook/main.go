package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned"
	"github.com/janekbaraniewski/kubeserial/pkg/images"
	"github.com/janekbaraniewski/kubeserial/pkg/webhooks"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.Install(scheme))
}

type HookParamters struct {
	certDir string
	port    int
}

func main() {
	var params HookParamters

	flag.IntVar(&params.port, "port", 8443, "Webhook port")
	flag.StringVar(&params.certDir, "certDir", "/certs/", "Webhook certificate folder")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	entryLog := setupLog.WithName("entrypoint")

	restConfig := ctrl.GetConfigOrDie()

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: ":8080",
		},
		HealthProbeBindAddress: ":8081",
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    params.port,
			CertDir: params.certDir,
		}),
	})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	clientset, err := versioned.NewForConfig(restConfig)
	if err != nil {
		entryLog.Error(err, "Failed to get clientset")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register(
		"/mutate-inject-device",
		&webhook.Admission{
			Handler: &webhooks.SerialDeviceInjector{
				Name:            "DeviceInjector",
				Clientset:       clientset,
				ConfigExtractor: images.NewOCIConfigExtractor(),
				Decoder:         admission.NewDecoder(scheme),
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
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
