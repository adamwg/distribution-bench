package tests

import (
	"github.com/adamwg/distribution-bench/config"
	"github.com/bloodorangeio/reggie"
)

func init() {
	registerCreator(config.BaseTest, createBaseTest)
}

// BaseTest tests the base route.
type BaseTest struct{}

// Prepare ...
func (bt *BaseTest) Prepare(client *reggie.Client) error {
	return nil
}

// Run ...
func (bt *BaseTest) Run(client *reggie.Client) (bytes uint64, err error) {
	resp, err := client.Do(client.NewRequest(reggie.GET, "/v2"))
	return uint64(len(resp.Body())), err
}

// Cleanup ...
func (bt *BaseTest) Cleanup(client *reggie.Client) error {
	return nil
}

func createBaseTest(cfg config.TestConfig) (Runner, error) {
	return &BaseTest{}, nil
}
