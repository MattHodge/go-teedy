package restore

import (
	"fmt"
	"path/filepath"

	"github.com/1set/gut/yos"

	"github.com/MattHodge/go-teedy/teedy"
)

// ViewDocumentFiles returns any attachments for a specific document
func (c *Client) ViewDocumentFiles(documentId string) ([]*yos.FilePathInfo, error) {
	fileAttachmentPath := filepath.Join(c.rootDocumentBackupDirectory, documentId, c.documentAttachmentDirectory)

	if yos.Exist(fileAttachmentPath) {
		return yos.ListFile(fileAttachmentPath)
	}

	// having no files directory is not an error, as some documents do not have files attached
	return nil, nil
}

func (c *Client) ViewDocuments() ([]*teedy.Document, error) {
	files, err := searchDirectoryForFiles(c.rootDocumentBackupDirectory, c.documentJSONFileBaseName)

	if err != nil {
		return nil, err
	}

	var res []*teedy.Document

	for _, file := range files {
		doc := new(teedy.Document)
		err := loadBackupFile(file, doc)

		if err != nil {
			return nil, fmt.Errorf("unable to load backup file: %w", err)
		}

		res = append(res, doc)
	}

	return res, nil
}

func (c *Client) Documents() error {
	// load docs from disk
	docs, err := c.ViewDocuments()

	if err != nil {
		return fmt.Errorf("cant load docs for restore: %w", err)
	}

	// load all the tags before restoring documents
	tags, err := c.client.Tag.GetAll()

	if err != nil {
		return fmt.Errorf("cant load tags for restore: %w", err)
	}

	for _, d := range docs {
		// check if doc exists
		existingDoc, err := c.client.Document.GetByTitle(d.Title)

		if existingDoc != nil {
			// delete it incase there are changes.. could be an update
			_, err := c.client.Document.Delete(existingDoc.Id)

			if err != nil {
				return fmt.Errorf("deleting existing document '%s' before storing failed: %w", d.Title, err)
			}
		}

		d.UpdateTagIDsByName(tags)

		_, err = c.client.Document.Add(d)

		if err != nil {
			return fmt.Errorf("cannot restore document %s: %w", d.Title, err)
		}
	}

	return nil
}
