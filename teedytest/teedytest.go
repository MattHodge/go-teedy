package teedytest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/teedy"

	"github.com/jarcoal/httpmock"
)

const (
	TeedyUsername              = "admin"
	TeedyPassword              = "admin"
	SkippingIntegrationMessage = "skipping integration test"
)

func SetupClient(t *testing.T) *teedy.Client {
	teedyUrl := os.Getenv("TEEDY_URL")

	if teedyUrl == "" {
		t.Fatalf("cannot run test because TEEDY_URL environment variable not set")
	}

	client, err := teedy.NewClient(teedyUrl, TeedyUsername, TeedyPassword)

	if err != nil {
		t.Fatalf("cannot run test because unable to get a new client: %v", err)
	}

	return client
}

// NewJsonResponder returns a JSON responder for httpmock.
func NewJsonResponder(s int, c string) httpmock.Responder {
	resp := httpmock.NewStringResponse(s, c)
	resp.Header.Set("Content-Type", "application/json")
	return httpmock.ResponderFromResponse(resp)
}

func LoadFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		t.Skipf("skipping test because unable unable to open file '%s': %v", path, err)
	}

	return file
}

func GetFileContents(t *testing.T, file string) []byte {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Skipf("skipping test because unable unable to read file: %s", err)
	}

	return b
}

func WriteToFile(t *testing.T, path string, content string) {
	basePath := filepath.Dir(path)
	err := os.MkdirAll(basePath, 0700)

	if err != nil {
		t.Skipf("cannot create directory %s: %s", basePath, err)
	}

	err = ioutil.WriteFile(path, []byte(content), 0644)

	if err != nil {
		t.Skipf("cannot write file %s: %s", path, err)
	}
}
