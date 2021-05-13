package backup_test

import (
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/teedytest"
	"github.com/jarcoal/httpmock"

	"github.com/MattHodge/go-teedy/backup"
	"github.com/MattHodge/go-teedy/teedy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {
	docFixture := `
{
 "total": 2,
 "documents": [
   {
     "id": "557aeee4-aa21-4369-8c1c-9705c842c673",
     "highlight": null,
     "file_id": "3e9ea4c9-e34b-4ec7-a2dc-083bd669f52f",
     "title": "Insurance",
     "description": "<p><ul><li>foo</li><li>bar</li></ul></p>",
     "create_date": 1615849200000,
     "update_date": 1618911359554,
     "language": "eng",
     "shared": false,
     "active_route": false,
     "current_step_name": null,
     "file_count": 5,
     "tags": [
       {
         "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
         "name": "baz",
         "color": "#3a87ad"
       }
     ]
   },
   {
     "id": "c819adab-a297-4a83-8816-a19ce6d82d20",
     "highlight": null,
     "file_id": null,
     "title": "Blinds",
     "description": "<ul><li>fii</li><li>bii</li><li>bee</li></ul>",
     "create_date": 1615417200000,
     "update_date": 1618352722495,
     "language": "eng",
     "shared": false,
     "active_route": false,
     "current_step_name": null,
     "file_count": 0,
     "tags": []
   }
 ],
 "suggestions": []
}
`
	docResponder := teedytest.NewJsonResponder(200, docFixture)
	httpmock.RegisterResponder("GET", "http://fake/api/document/list", docResponder)

	fileFixture := `
{
 "files": [
   {
     "id": "1",
     "name": "file1.png",
     "document_id": "557aeee4-aa21-4369-8c1c-9705c842c673"
   },
   {
     "id": "2",
     "name": "file2.png",
     "document_id": "557aeee4-aa21-4369-8c1c-9705c842c673"
   }
 ]
}
`
	fileResponder := teedytest.NewJsonResponder(200, fileFixture)
	httpmock.RegisterResponder("GET", "http://fake/api/file/list", fileResponder)

	file1DataResponder := httpmock.NewBytesResponder(200, teedytest.GetFileContents(t, "testdata/image.png"))
	httpmock.RegisterResponder("GET", "http://fake/api/file/1/data", file1DataResponder)

	file2DataResponder := httpmock.NewBytesResponder(200, teedytest.GetFileContents(t, "testdata/image.png"))
	httpmock.RegisterResponder("GET", "http://fake/api/file/2/data", file2DataResponder)

	tc := teedy.NewFakeClient()
	backupDir := t.TempDir()
	backupClient := backup.NewBackupClient(tc, backupDir)

	err := backupClient.Documents()

	require.NoError(t, err)
	assert.FileExists(t, backupClient.DocumentBackupJSONFilePath("557aeee4-aa21-4369-8c1c-9705c842c673")) //backupdir/documents/<guid>/document.json
	assert.FileExists(t, backupClient.DocumentBackupJSONFilePath("c819adab-a297-4a83-8816-a19ce6d82d20"))
	assert.DirExists(t, backupClient.DocumentBackupDirectory("557aeee4-aa21-4369-8c1c-9705c842c673")) //backupdir/documents/<guid>/files
	assert.FileExists(t, filepath.Join(backupClient.DocumentBackupDirectory("557aeee4-aa21-4369-8c1c-9705c842c673"), "files", "file1.png"))
	assert.FileExists(t, filepath.Join(backupClient.DocumentBackupDirectory("557aeee4-aa21-4369-8c1c-9705c842c673"), "files", "file2.png"))
}
