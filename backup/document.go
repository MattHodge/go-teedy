package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MattHodge/go-teedy/teedy"
)

type DocumentBackup struct {
	FullDirectory        string
	FullDirectoryFiles   string
	FullPathDocumentJSON string
	Document             *teedy.Document
}

func (s *DocumentBackup) Save() error {
	os.MkdirAll(s.FullDirectory, 0700)
	os.MkdirAll(s.FullDirectoryFiles, 0700)
	err := dumpJson(s.Document, s.FullPathDocumentJSON)

	if err != nil {
		return fmt.Errorf("cannot save: %w", err)
	}

	return nil
}

func Document(document *teedy.Document, basePath string) *DocumentBackup {
	fullDirectory := filepath.Join(basePath, "documents", document.Id)
	return &DocumentBackup{
		FullDirectory:        fullDirectory,
		FullDirectoryFiles:   filepath.Join(fullDirectory, "files"),
		FullPathDocumentJSON: filepath.Join(fullDirectory, DOCUMENT_BACKUP_FILENAME),
		Document:             document,
	}
}
