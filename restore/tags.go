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
	tagsFromDisk, err := c.ViewTags()

	if err != nil {
		return fmt.Errorf("cant load tagsFromDisk for restore: %w", err)
	}

	restoredTags := make(map[string]*teedy.Tag)
	var notRestoredYet []*teedy.Tag

	for _, tagFromDisk := range tagsFromDisk {
		err := c.deleteTagIfExistsInDestination(tagFromDisk)

		if err != nil {
			return err
		}

		// first restore tags with no parents
		if tagFromDisk.Parent == "" {
			restoredTag, err := c.client.Tag.Add(tagFromDisk)

			if err != nil {
				return err
			}

			// create a map so we can us the old tag id to find the new tag
			restoredTags[tagFromDisk.Id] = restoredTag
		} else {
			notRestoredYet = append(notRestoredYet, tagFromDisk)
		}
	}

	// restore the rest of the tags
	for _, tagFromDisk := range notRestoredYet {
		// try find its parent by name
		if parentTag, ok := restoredTags[tagFromDisk.Parent]; ok {
			tagFromDisk.Parent = parentTag.Id
		}

		restoredTag, err := c.client.Tag.Add(tagFromDisk)

		if err != nil {
			return err
		}

		restoredTags[tagFromDisk.Id] = restoredTag
	}

	return nil
}

func (c *Client) deleteTagIfExistsInDestination(t *teedy.Tag) error {
	// check if tag exists
	existingTag, err := c.client.Tag.GetByName(t.Name)

	if err != nil {
		return err
	}

	if existingTag != nil {
		// delete it incase there are changes.. could be an update
		_, err := c.client.Tag.Delete(existingTag.Id)

		if err != nil {
			return err
		}
	}

	return nil
}
