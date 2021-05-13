package teedy_test

import (
	"testing"

	"github.com/MattHodge/go-teedy/teedy"
)

const (
	testTeedyURL                = "http://localhost:8080"
	testTeedyUsername           = "admin"
	testTeedyPassword           = "superSecure"
	testSkippingIntegrationTest = "skipping integration test"
)

func setup(t *testing.T) *teedy.Client {
	client, err := teedy.NewClient(testTeedyURL, testTeedyUsername, testTeedyPassword)

	if err != nil {
		t.Skipf("skipping test because unable to get a new client")
	}

	return client
}
