package tests

import "github.com/bloodorangeio/reggie"

// Runner runs a test trial.
type Runner interface {
	// Prepare does any necessary setup for the test, such as pushing a blob to
	// be fetched as part of the test.
	Prepare(client *reggie.Client) error
	// Run runs the test. It returns the number of bytes transmitted or received
	// (depending on the test type) and any error encountered.
	Run(client *reggie.Client) (bytes uint64, err error)
	// Cleanup does any necessary cleanup post-test, such as deleting a pushed
	// blob.
	Cleanup(client *reggie.Client) error
}
