package evernote_test

import (
	"testing"

	"github.com/MattHodge/go-teedy/evernote"
	"github.com/MattHodge/go-teedy/teedytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImport(t *testing.T) {
	client := teedytest.SetupClient(t)
	ec := evernote.NewImportClient("testdata/test.enex", client)
	res, err := ec.Import()

	require.NoError(t, err)
	require.Len(t, res, 2, "imported amount of evernote notes not expected")
	require.NotNil(t, res[0].Id, "created document should have id")
	require.NotNil(t, res[1].Id, "created document should have id")

	doc, err := client.Document.Get(res[0].Id)
	require.NoError(t, err, "getting created document from server should not error")

	assert.Equal(t, "1. test", doc.Title, "document title should match imported")
	assert.Contains(t, doc.Description, "some text", "document description should match imported")

	files, err := client.File.Get(res[0].Id)
	require.NoError(t, err, "getting files for created document should not error")

	require.Len(t, files, 2, "should be 2 attachments for imported evernote note")
}
