package images

import (
	"context"

	"github.com/regclient/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
	v1 "github.com/regclient/regclient/types/oci/v1"
	"github.com/regclient/regclient/types/ref"
)

type OCIConfigExtractor struct {
	client regclient.RegClient
}

func (e *OCIConfigExtractor) getManifestFromList(ctx context.Context, man manifest.Manifest, r ref.Ref) (manifest.Manifest, ref.Ref, error) {
	manifests, err := man.GetManifestList()

	if err != nil {
		return nil, r, err
	}

	r.Digest = string(manifests[0].Digest)
	specMan, err := e.client.ManifestGet(ctx, r)

	if err != nil {
		return nil, r, err
	}

	return specMan, r, nil
}

func (e *OCIConfigExtractor) GetImageConfig(ctx context.Context, image string) (conf v1.ImageConfig, err error) {
	referance, err := ref.New(image)

	if err != nil {
		return v1.ImageConfig{}, err
	}

	manifest, err := e.client.ManifestGet(ctx, referance)

	if err != nil {
		return v1.ImageConfig{}, err
	}

	if manifest.IsList() {
		man, ref, err := e.getManifestFromList(ctx, manifest, referance)
		if err != nil {
			return v1.ImageConfig{}, err
		}
		manifest = man
		referance = ref
	}

	config, err := manifest.GetConfig()
	if err != nil {
		return v1.ImageConfig{}, err
	}

	blobConfig, err := e.client.BlobGetOCIConfig(ctx, referance, config)
	if err != nil {
		return v1.ImageConfig{}, err
	}
	return blobConfig.GetConfig().Config, nil
}

func NewOCIConfigExtractor() *OCIConfigExtractor {
	return &OCIConfigExtractor{
		client: regclient.NewRegClient(),
	}
}
