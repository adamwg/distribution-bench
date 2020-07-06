package config

import "encoding/json"

// Config is configuration for the distribution benchmarking tool.
type Config struct {
	// Registry configures the registry to connect to.
	Registry RegistryConfig `json:"registry"`
	// Tests configures tests to run.
	Tests []TestConfig `json:"tests"`

	// Debug enables debug output.
	Debug bool `json:"debug"`
}

// RegistryConfig configures access to a registry for testing.
type RegistryConfig struct {
	// RootURL is the root URL of the registry. This URL should *not* include
	// the `/v2` prefix used by the distribution API.
	RootURL string `json:"root_url"`
	// Username is a username for basic authentication. If empty, authentication
	// will not be used.
	Username string `json:"username"`
	// Password is a password for basic authentication.
	Password string `json:"password"`
	// Namespace is a namespace to use within the registry. You may need to
	// pre-create this namespace in some registry implementations. If empty, a
	// random namespace with two path components will be generated.
	Namespace string `json:"namespace"`
}

// TestConfig configures a single test to run.
type TestConfig struct {
	// Type is the type of test.
	Type TestType `json:"type"`
	// Trials is the number of total trials to run. Trials may be run
	// concurrently, depending on the Concurrency setting.
	Trials int `json:"trials"`
	// Concurrency is how many tests to run concurrently.
	Concurrency int `json:"concurrency"`
	// Parameters contains test-specific parameters.
	Parameters json.RawMessage `json:"parameters"`
}

// TestType is a type of test to run.
type TestType string

const (
	// BaseTest is a test for the base route, and has no practical utility.
	BaseTest TestType = "base"

	// GetManifestTest is a test type for downloading manifests.
	GetManifestTest TestType = "get-manifest"
	// PutManifestTest is a test type for uploading manifests.
	PutManifestTest TestType = "put-manifest"

	// GetBlobTest is a test type for downloading blobs.
	GetBlobTest TestType = "get-blob"
	// PutBlobMonolithicTest is a test type for uploading blobs using the
	// monolithic (PUT) method.
	PutBlobMonolithicTest TestType = "put-blob-monolithic"
	// PutBlobStreamingTest is a test type for uplaoding blobs using the
	// streaming (PATCH) method.
	PutBlobStreamingTest TestType = "put-blob-streaming"
)
