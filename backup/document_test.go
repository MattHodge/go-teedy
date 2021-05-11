package backup_test

import (
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/backup"
	"github.com/MattHodge/go-teedy/teedy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDocumentBackup(t *testing.T) {
	baseDir := t.TempDir()
	document := &teedy.Document{
		Title: "abc",
		Id:    "123",
	}
	got := backup.Document(document, baseDir)
	want := &backup.DocumentBackup{
		FullDirectory:        filepath.Join(baseDir, "documents", "123"),
		FullDirectoryFiles:   filepath.Join(baseDir, "documents", "123", "files"),
		FullPathDocumentJSON: filepath.Join(baseDir, "documents", "123", backup.DOCUMENT_BACKUP_FILENAME),
		Document:             document,
	}

	assert.EqualValues(t, want, got)
}

func TestGetDocumentBackup_WithSave(t *testing.T) {
	baseDir := t.TempDir()
	document := &teedy.Document{
		Title: "abc",
		Id:    "123",
	}

	backupDoc := backup.Document(document, baseDir)
	err := backupDoc.Save()

	require.NoError(t, err)
	assert.DirExists(t, backupDoc.FullDirectory)
	assert.DirExists(t, backupDoc.FullDirectoryFiles)
	assert.FileExists(t, backupDoc.FullPathDocumentJSON)
}
