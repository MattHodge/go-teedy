package teedy_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/MattHodge/go-teedy/teedy"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileService_Get(t *testing.T) {
	fixture := `
{
    "files": [
        {
            "id": "3e9ea4c9-e34b-4ec7-a2dc-083bd669f52f",
            "processing": false,
            "name": "Foo.pdf",
            "version": 0,
            "mimetype": "application/pdf",
            "document_id": "557aeee4-aa21-4369-8c1c-9705c842c673",
            "create_date": 1618031146082,
            "size": 182530
        },
        {
            "id": "87e442a5-1355-499f-921c-44f1d1e1d9db",
            "processing": false,
            "name": "Bar.pdf",
            "version": 0,
            "mimetype": "application/pdf",
            "document_id": "557aeee4-aa21-4369-8c1c-9705c842c673",
            "create_date": 1618031147088,
            "size": 172194
        }
    ]
}
`
	responder := newJsonResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://fake/api/file/list?id=1", responder)
	client := teedy.NewFakeClient()

	files, err := client.File.Get("1")

	require.NoError(t, err)
	assert.Len(t, files, 2)
	assert.Equal(t, files[0].Name, "Foo.pdf")
}

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

	file := loadFile(t, "testdata/image.png")
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

	file := loadFile(t, "testdata/image.png")
	defer file.Close()

	// test adding the file
	s, err := client.File.Add(doc.Id, "", file)
	require.NoError(t, err)
	assert.Equal(t, s.Status, "ok")

	// test the document has the file attached
	readDoc := readTestDocument(t, client, doc.Id)
	assert.Equal(t, 1, readDoc.FileCount)
}

func TestFileService_GetData_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)
	doc := createTestDocument(t, client)

	// add file to it
	file := loadFile(t, "testdata/image.png")
	defer file.Close()
	addedFile, err := client.File.Add(doc.Id, "", file)
	require.NoError(t, err, "should not error adding file")

	got, err := client.File.GetData(addedFile.Id)

	require.NoError(t, err, "should be no error getting zipped file")

	assert.IsType(t, got, []byte{})

	uploadedFileContents := getFileContents(t, "testdata/image.png")
	assert.Equal(t, uploadedFileContents, got, "uploaded and read file should match")
}

func TestFileService_GetZippedFiles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)
	doc := createTestDocument(t, client)

	// add file to it
	file := loadFile(t, "testdata/image.png")
	defer file.Close()
	_, err := client.File.Add(doc.Id, "", file)

	got, err := client.File.GetZippedFiles(doc.Id)

	require.NoError(t, err, "should be no error getting zipped file")

	var want []byte
	assert.IsType(t, got, want)
}

func loadFile(t *testing.T, path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		t.Skipf("skipping test because unable unable to open file '%s': %v", path, err)
	}

	return file
}

func getFileContents(t *testing.T, file string) []byte {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		t.Skipf("skipping test because unable unable to read file: %s", err)
	}

	return b
}

// createTestDocument creates a document in the teedy API for integration test purposes
func createTestDocument(t *testing.T, client *teedy.Client) *teedy.Document {
	doc := teedy.NewDocument("test document", "eng")
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
