package teedy_test

import (
	"os"
	"testing"

	"github.com/MattHodge/go-teedy/teedy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileService_GetAll_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)

	_, err := client.File.GetAll()

	require.NoError(t, err, "should be no error getting zipped file")
}

func TestFileService_Add_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}
	client := setup(t)

	file, err := loadFile(t, "testdata/image.png")
	defer file.Close()
	s, err := client.File.Add("", "", file)
	require.NoError(t, err)
	assert.Equal(t, s.Status, "ok")
}

func TestFileService_AddToDocument_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)

	doc := createTestDocument(t, client)

	file, err := loadFile(t, "testdata/image.png")
	defer file.Close()

	// test adding the file
	s, err := client.File.Add(doc.Id, "", file)
	require.NoError(t, err)
	assert.Equal(t, s.Status, "ok")

	// test the document has the file attached
	readDoc := readTestDocument(t, client, doc.Id)
	assert.Equal(t, 1, readDoc.FileCount)
}

func TestFileService_GetZippedFiles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)
	doc := createTestDocument(t, client)

	// add file to it
	file, err := loadFile(t, "testdata/image.png")
	defer file.Close()
	_, err = client.File.Add(doc.Id, "", file)

	got, err := client.File.GetZippedFiles(doc.Id)

	require.NoError(t, err, "should be no error getting zipped file")

	var want []byte
	assert.IsType(t, got, want)
}

func loadFile(t *testing.T, path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		t.Skipf("skipping test because unable unable to open file '%s': %v", path, err)
	}

	return file, nil
}

// createTestDocument creates a document in the teedy API for integration test purposes
func createTestDocument(t *testing.T, client *teedy.Client) *teedy.Document {
	doc, err := teedy.NewDocument("test document", "eng")
	d, err := client.Document.Add(doc)

	if err != nil {
		t.Skipf("skipping test due to error: %v", err)
	}

	return d
}

func readTestDocument(t *testing.T, client *teedy.Client, id string) *teedy.Document {
	doc, err := client.Document.Get(id)

	if err != nil {
		t.Skipf("skipping test due to error: %v", err)
	}

	return doc
}
