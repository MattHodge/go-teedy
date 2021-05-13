package backup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/MattHodge/go-teedy/teedy"
)

type Client struct {
	client                      *teedy.Client
	rootBackupDirectory         string
	rootTagBackupDirectory      string
	rootDocumentBackupDirectory string
	tagJSONFileBaseName         string
	documentJSONFileBaseName    string
}

func NewBackupClient(client *teedy.Client, directory string) *Client {
	return &Client{
		client:                      client,
		rootBackupDirectory:         directory,
		rootTagBackupDirectory:      filepath.Join(directory, "tags"),
		rootDocumentBackupDirectory: filepath.Join(directory, "documents"),
		tagJSONFileBaseName:         "tag.json",
		documentJSONFileBaseName:    "document.json",
	}
}

func dumpJson(i interface{}, path string) error {
	bytes, err := json.MarshalIndent(i, "", "    ")

	if err != nil {
		return fmt.Errorf("error marshalling to json: %v", err)
	}

	return ioutil.WriteFile(path, bytes, 0644)
}
