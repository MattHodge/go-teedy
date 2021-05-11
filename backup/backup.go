package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MattHodge/go-teedy/teedy"
)

const (
	TAG_BACKUP_FILENAME      = "tag.json"
	DOCUMENT_BACKUP_FILENAME = "document.json"
)

type BackupClient struct {
	client          *teedy.Client
	backupDirectory string
}

func dumpJson(i interface{}, path string) error {
	bytes, err := json.MarshalIndent(i, "", "    ")

	if err != nil {
		return fmt.Errorf("error marshalling to json: %v", err)
	}

	return ioutil.WriteFile(path, bytes, 0644)
}
