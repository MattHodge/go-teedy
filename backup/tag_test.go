package backup_test

import (
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/teedytest"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/MattHodge/go-teedy/backup"

	"github.com/MattHodge/go-teedy/teedy"
	"github.com/stretchr/testify/assert"
)

func TestTagBackupJSONFilePath(t *testing.T) {
	tc := teedy.NewFakeClient()
	backupDir := t.TempDir()
	backupClient := backup.NewBackupClient(tc, backupDir)

	got := backupClient.TagBackupJSONFilePath("123")
	want := filepath.Join(backupDir, "tags", "123", "tag.json")

	assert.Equal(t, want, got)
}

func TestTagBackupDirectory(t *testing.T) {
	tc := teedy.NewFakeClient()
	backupDir := t.TempDir()
	backupClient := backup.NewBackupClient(tc, backupDir)

	got := backupClient.TagBackupDirectory("123")
	want := filepath.Join(backupDir, "tags", "123")

	assert.Equal(t, want, got)
}

func TestTag(t *testing.T) {
	fixture := `
{
 "tags": [
   {
     "id": "d12a2911-c70f-4498-9c2b-19ba2646ca84",
     "name": "tax",
     "color": "#3a87ad"
   },
   {
     "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
     "name": "tax_2020",
     "color": "#3a87ad",
     "parent": "d12a2911-c70f-4498-9c2b-19ba2646ca84"
   }
 ]
}
`
	responder := teedytest.NewJsonResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://fake/api/tag/list", responder)

	tc := teedy.NewFakeClient()
	backupDir := t.TempDir()
	backupClient := backup.NewBackupClient(tc, backupDir)

	err := backupClient.Tags()

	require.NoError(t, err)
	assert.FileExists(t, backupClient.TagBackupJSONFilePath("d12a2911-c70f-4498-9c2b-19ba2646ca84"))
	assert.FileExists(t, backupClient.TagBackupJSONFilePath("bcd8e09b-84bc-4926-afce-222b7c21d8eb"))
}
