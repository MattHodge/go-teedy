package backup

import (
	"fmt"
	"os"
	"path/filepath"
)

func (b *Client) TagBackupJSONFilePath(tagId string) string {
	return filepath.Join(b.TagBackupDirectory(tagId), b.tagJSONFileBaseName)
}

func (b *Client) TagBackupDirectory(tagId string) string {
	return filepath.Join(b.rootTagBackupDirectory, tagId)
}

func (b *Client) Tags() error {
	tags, err := b.client.Tag.GetAll()

	if err != nil {
		return fmt.Errorf("cannot get tags: %w", err)
	}

	for _, tag := range tags {
		os.MkdirAll(b.TagBackupDirectory(tag.Id), 0700)
		err := dumpJson(tag, b.TagBackupJSONFilePath(tag.Id))

		if err != nil {
			return fmt.Errorf("cannot save: %w", err)
		}
	}

	return nil
}
