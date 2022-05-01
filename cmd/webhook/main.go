package main

import (
	"crypto/sha256"
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/janekbaraniewski/kubeserial/pkg/webhooks"
)

var setupLog = ctrl.Log.WithName("setup")

type HookParamters struct {
	certDir       string
	sidecarConfig string
	port          int
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		*files = append(*files, path)
		return nil
	}
}

func loadConfig(configFile string) (*webhooks.Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	setupLog.Info("New configuration", "sha256sum", sha256.Sum256(data))

	var cfg webhooks.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	var params HookParamters

	flag.IntVar(&params.port, "port", 8443, "Wehbook port")
	flag.StringVar(&params.certDir, "certDir", "/certs/", "Wehbook certificate folder")
	flag.StringVar(&params.sidecarConfig, "sidecarConfig", "/etc/webhook/config/sidecarconfig.yaml", "Wehbook sidecar config")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	entryLog := setupLog.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	config, err := loadConfig(params.sidecarConfig)

	if err != nil {
		entryLog.Error(err, "Can't load config")
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	hookServer.Port = params.port
	hookServer.CertDir = params.certDir

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/mutate", &webhook.Admission{Handler: &webhooks.SidecarInjector{Name: "Logger", Client: mgr.GetClient(), SidecarConfig: config}})

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
