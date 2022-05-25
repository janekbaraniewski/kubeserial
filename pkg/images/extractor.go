package images

import (
	"context"

	"github.com/regclient/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
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

func (e *OCIConfigExtractor) GetEntrypointAndCMD(ctx context.Context, image string) (entrypoint []string, cmd []string, err error) {
	referance, err := ref.New(image)

	if err != nil {
		panic(err)
	}

	manifest, err := e.client.ManifestGet(context.Background(), referance)

	if manifest.IsList() {
		man, ref, err := e.getManifestFromList(ctx, manifest, referance)
		if err != nil {
			return nil, nil, err
		}
		manifest = man
		referance = ref
	}

	config, err := manifest.GetConfig()
	if err != nil {
		return nil, nil, err
	}

	blobConfig, err := e.client.BlobGetOCIConfig(ctx, referance, config)
	if err != nil {
		return nil, nil, err
	}
	conf := blobConfig.GetConfig().Config
	return conf.Entrypoint, conf.Cmd, nil
}

func NewOCIConfigExtractor() *OCIConfigExtractor {
	return &OCIConfigExtractor{
		client: regclient.NewRegClient(),
	}
}
