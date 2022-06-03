package utils

import (
	"k8s.io/apimachinery/pkg/util/yaml"
)

func LoadResourceFromYaml(fs FileSystem, filepath string, data interface{}) error {
	reader, err := fs.Open(filepath)
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
