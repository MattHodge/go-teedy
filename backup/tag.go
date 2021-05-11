package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MattHodge/go-teedy/teedy"
)

type TagBackup struct {
	FullDirectory        string
	FullPathDocumentJSON string
	Tag                  *teedy.Tag
}

func Tag(tag *teedy.Tag, basePath string) *TagBackup {
	fullDirectory := filepath.Join(basePath, "tags", tag.Id)
	return &TagBackup{
		FullDirectory:        fullDirectory,
		FullPathDocumentJSON: filepath.Join(fullDirectory, TAG_BACKUP_FILENAME),
		Tag:                  tag,
	}
}

func (t *TagBackup) Save() error {
	os.MkdirAll(t.FullDirectory, 0700)
	err := dumpJson(t.Tag, t.FullPathDocumentJSON)

	if err != nil {
		return fmt.Errorf("cannot save: %w", err)
	}

	return nil
}
