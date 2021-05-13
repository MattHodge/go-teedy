package backup

import (
	"fmt"
	"os"
	"path/filepath"
)

func (b *Client) DocumentBackupJSONFilePath(documentId string) string {
	return filepath.Join(b.DocumentBackupDirectory(documentId), b.documentJSONFileBaseName)
}

func (b *Client) DocumentBackupDirectory(documentId string) string {
	return filepath.Join(b.rootDocumentBackupDirectory, documentId)
}

func (b *Client) Documents() error {
	docs, err := b.client.Document.GetAll()

	if err != nil {
		return fmt.Errorf("cannot get tags: %w", err)
	}

	for _, doc := range docs.Documents {
		os.MkdirAll(b.DocumentBackupDirectory(doc.Id), 0700)
		err := dumpJson(doc, b.DocumentBackupJSONFilePath(doc.Id))

		if err != nil {
			return fmt.Errorf("cannot save: %w", err)
		}

		// get files
		docFiles, err := b.client.File.Get(doc.Id)

		if err != nil {
			return fmt.Errorf("cannot get files for doc id %s: %w", doc.Id, err)
		}

		for _, file := range docFiles {
			err := b.File(file)
			if err != nil {
				return fmt.Errorf("cannot save file id %s: %w", file.Id, err)
			}
		}
	}

	return nil
}
