package restore_test

import (
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/teedytest"

	"github.com/MattHodge/go-teedy/restore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBackupDocuments(t *testing.T) {
	baseDir := t.TempDir()
	doc1 := `
{
    "id": "1",
    "file_id": "3e9ea4c9-e34b-4ec7-a2dc-083bd669f52f",
    "title": "foo",
    "tags": [
        {
            "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
            "name": "tax_2020",
            "color": "#3a87ad"
        }
    ]
}
`

	doc2 := `
{
    "id": "2",
    "file_id": "3e9ea4c9-e34b-4ec7-a2dc-083bd669f52f",
    "title": "bar",
    "tags": [
        {
            "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
            "name": "tax_2020",
            "color": "#3a87ad"
        }
    ]
}
`

	teedytest.WriteToFile(t, filepath.Join(baseDir, "documents", "1", "document.json"), doc1)
	teedytest.WriteToFile(t, filepath.Join(baseDir, "documents", "2", "document.json"), doc2)

	docs, err := restore.LoadBackupDocuments(filepath.Join(baseDir, "documents"))

	require.NoError(t, err)
	require.Len(t, docs, 2)
	assert.Equal(t, "1", docs[0].Id)
	assert.Equal(t, "2", docs[1].Id)
}
