package utils

import (
	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func LoadResourceFromYaml(fs FileSystem, filepath kubeserial.ResourceSpecPath, data interface{}) error {
	reader, err := fs.Open(string(filepath))
	if err != nil {
		return err
	}
	defer reader.Close()
	err = yaml.NewYAMLOrJSONDecoder(reader, 4096).Decode(data)
	if err != nil {
		return err
	}

	return nil
}
