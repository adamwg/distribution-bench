package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/adamwg/distribution-bench/config"
	"github.com/bloodorangeio/reggie"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

func init() {
	registerCreator(config.PutBlobMonolithicTest, createPutBlobMonolithicTest)
	registerCreator(config.PutBlobStreamingTest, createPutBlobStreamingTest)
}

type putBlobStreamingTest struct {
	PutBlobTestParameters
}

// Prepare ...
func (pbt *putBlobStreamingTest) Prepare(client *reggie.Client) error {
	return nil
}

// Run ...
func (pbt *putBlobStreamingTest) Run(client *reggie.Client) (uint64, error) {
	req := client.NewRequest(reggie.POST, "/v2/<name>/blobs/uploads/")
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "initiating upload")
	}
	if code := resp.StatusCode(); code != http.StatusCreated && code != http.StatusAccepted {
		return 0, fmt.Errorf("initiating upload: expected created or accepted, got %d", code)
	}

	uploadPath := resp.GetRelativeLocation()

	body := digestingRandomReader(pbt.SizeBytes)
	resetBody := func(r *reggie.Request) error {
		body = digestingRandomReader(pbt.SizeBytes)
		r.SetBody(body)
		return nil
	}

	req = client.NewRequest(reggie.PATCH, uploadPath, reggie.WithRetryCallback(resetBody)).
		SetHeader("Content-Type", mediaTypeOctetStream).
		SetBody(body)
	resp, err = client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "uploading data")
	}
	if code := resp.StatusCode(); code != http.StatusAccepted {
		return 0, fmt.Errorf("uploading data: expected accepted, got %d", code)
	}

	dg := body.Digest()
	req = client.NewRequest(reggie.PUT, resp.GetRelativeLocation()).
		SetQueryParam("digest", dg.String()).
		SetHeader("Content-Type", mediaTypeOctetStream).
		SetHeader("Content-Length", strconv.FormatUint(pbt.SizeBytes, 10))
	resp, err = client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "uploading data")
	}
	if code := resp.StatusCode(); code != http.StatusCreated {
		return 0, fmt.Errorf("uploading data: expected created, got %d", code)
	}

	return pbt.SizeBytes, nil
}

// Cleanup ...
func (pbt *putBlobStreamingTest) Cleanup(client *reggie.Client) error {
	return nil
}

type putBlobMonolithicTest struct {
	body   []byte
	digest digest.Digest
}

// Prepare ...
func (pbt *putBlobMonolithicTest) Prepare(client *reggie.Client) error {
	return nil
}

// Run ...
func (pbt *putBlobMonolithicTest) Run(client *reggie.Client) (bytes uint64, err error) {
	req := client.NewRequest(reggie.POST, "/v2/<name>/blobs/uploads/").
		SetHeader("Content-Length", strconv.Itoa(len(pbt.body))).
		SetHeader("Content-Type", mediaTypeOctetStream).
		SetQueryParam("digest", pbt.digest.String()).
		SetBody(pbt.body)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if code := resp.StatusCode(); code != http.StatusCreated && code != http.StatusAccepted {
		return 0, fmt.Errorf("uploading blob: expected created or accepted, got %d", code)
	}

	return uint64(len(pbt.body)), nil
}

// Cleanup ...
func (pbt *putBlobMonolithicTest) Cleanup(client *reggie.Client) error {
	return nil
}

// PutBlobTestParameters is the parameters for a put blob test.
type PutBlobTestParameters struct {
	// SizeBytes is the total size of the blob upload, in bytes.
	SizeBytes uint64 `json:"size_bytes"`
}

func createPutBlobMonolithicTest(cfg config.TestConfig) (Runner, error) {
	var params PutBlobTestParameters
	err := json.Unmarshal(cfg.Parameters, &params)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling test parameters")
	}

	body, err := randomBytes(int64(params.SizeBytes))
	if err != nil {
		return nil, errors.Wrap(err, "creating random blob")
	}
	digest := digest.Canonical.FromBytes(body)

	return &putBlobMonolithicTest{
		body:   body,
		digest: digest,
	}, nil
}

func createPutBlobStreamingTest(cfg config.TestConfig) (Runner, error) {
	var params PutBlobTestParameters
	err := json.Unmarshal(cfg.Parameters, &params)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling test parameters")
	}

	return &putBlobStreamingTest{
		PutBlobTestParameters: params,
	}, nil
}

type digestingReader struct {
	io.Reader
	digest.Digester
}

func digestingRandomReader(length uint64) *digestingReader {
	underlying := randomReader(int64(length))
	digester := digest.Canonical.Digester()
	tee := io.TeeReader(underlying, digester.Hash())

	return &digestingReader{
		Reader:   tee,
		Digester: digester,
	}
}
