package restore_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/backup"
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

	writeToFile(t, filepath.Join(baseDir, "documents", "1", backup.DOCUMENT_BACKUP_FILENAME), doc1)
	writeToFile(t, filepath.Join(baseDir, "documents", "2", backup.DOCUMENT_BACKUP_FILENAME), doc2)

	docs, err := restore.LoadBackupDocuments(filepath.Join(baseDir, "documents"))

	require.NoError(t, err)
	require.Len(t, docs, 2)
	assert.Equal(t, "1", docs[0].Id)
	assert.Equal(t, "2", docs[1].Id)
}

func writeToFile(t *testing.T, path string, content string) {
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
