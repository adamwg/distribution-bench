package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adamwg/distribution-bench/config"
	"github.com/bloodorangeio/reggie"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	imagespec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

func init() {
	registerCreator(config.PutManifestTest, createPutManifestTest)
}

type putManifestTest struct {
	reference       string
	configDigest    string
	configContent   []byte
	manifestContent []byte
}

// Prepare ...
func (pmt *putManifestTest) Prepare(client *reggie.Client) error {
	return nil
}

// Run ...
func (pmt *putManifestTest) Run(client *reggie.Client) (bytes uint64, err error) {
	configReq := client.NewRequest(reggie.POST, "/v2/<name>/blobs/uploads/").
		SetHeader("Content-Length", strconv.Itoa(len(pmt.configContent))).
		SetHeader("Content-Type", mediaTypeOctetStream).
		SetQueryParam("digest", pmt.configDigest).
		SetBody(pmt.configContent)

	resp, err := client.Do(configReq)
	if err != nil {
		return 0, err
	}
	if code := resp.StatusCode(); code != http.StatusCreated && code != http.StatusAccepted {
		return 0, fmt.Errorf("uploading config: expected created or accepted, got %d", code)
	}

	manifestReq := client.NewRequest(
		reggie.PUT,
		"/v2/<name>/manifests/<reference>",
		reggie.WithReference(pmt.reference),
	).
		SetHeader("Content-Type", "application/vnd.oci.image.manifest.v1+json").
		SetBody(pmt.manifestContent)

	resp, err = client.Do(manifestReq)
	if err != nil {
		return 0, err
	}
	if code := resp.StatusCode(); code != http.StatusCreated {
		return 0, fmt.Errorf("uploading manifest: expected created, got %d", code)
	}

	return uint64(len(pmt.manifestContent) + len(pmt.configContent)), nil
}

// Cleanup ...
func (pmt *putManifestTest) Cleanup(client *reggie.Client) error {
	// TODO(awg): Delete the manifest we pushed.
	return nil
}

func createPutManifestTest(cfg config.TestConfig) (Runner, error) {
	config := imagespec.Image{
		Architecture: "amd64",
		OS:           "linux",
		RootFS: imagespec.RootFS{
			Type:    "layers",
			DiffIDs: []digest.Digest{},
		},
	}
	configBlobContent, err := json.MarshalIndent(&config, "", "\t")
	if err != nil {
		return nil, errors.Wrap(err, "marshaling image config")
	}

	configBlobDigest := digest.FromBytes(configBlobContent)

	manifest := imagespec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		Config: imagespec.Descriptor{
			MediaType: "application/vnd.oci.image.config.v1+json",
			Digest:    configBlobDigest,
			Size:      int64(len(configBlobContent)),
		},
		Layers: []imagespec.Descriptor{},
	}

	manifestContent, err := json.MarshalIndent(&manifest, "", "\t")
	if err != nil {
		return nil, errors.Wrap(err, "marshaling manifest")
	}

	return &putManifestTest{
		reference:       randomTag(),
		configDigest:    configBlobDigest.String(),
		configContent:   configBlobContent,
		manifestContent: manifestContent,
	}, nil
}

func randomTag() string {
	return randomString(10)
}
