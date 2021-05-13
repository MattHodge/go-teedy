package restore

import (
	"fmt"

	"github.com/MattHodge/go-teedy/teedy"
)

func (c *Client) ViewTags() ([]*teedy.Tag, error) {
	files, err := searchDirectoryForFiles(c.rootTagBackupDirectory, c.tagJSONFileBaseName)

	if err != nil {
		return nil, err
	}

	var res []*teedy.Tag

	for _, file := range files {
		tag := new(teedy.Tag)
		err := loadBackupFile(file, tag)

		if err != nil {
			return nil, fmt.Errorf("unable to load backup file: %w", err)
		}

		res = append(res, tag)
	}

	return res, nil
}

func (c *Client) Tags() error {
	// load docs from disk
	tags, err := c.ViewTags()

	if err != nil {
		return fmt.Errorf("cant load tags for restore: %w", err)
	}

	for _, t := range tags {
		// check if tag exists
		existingTag, err := c.client.Tag.GetByName(t.Name)

		if existingTag != nil {
			// delete it incase there are changes.. could be an update
			_, err := c.client.Tag.Delete(existingTag.Id)

			if err != nil {
				return err
			}
		}

		_, err = c.client.Tag.Add(t)

		if err != nil {
			return err
		}
	}

	return nil
}
