package restore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/1set/gut/yos"
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

func NewRestoreClient(client *teedy.Client, directory string) *Client {
	return &Client{
		client:                      client,
		rootBackupDirectory:         directory,
		rootTagBackupDirectory:      filepath.Join(directory, "tags"),
		rootDocumentBackupDirectory: filepath.Join(directory, "documents"),
		tagJSONFileBaseName:         "tag.json",
		documentJSONFileBaseName:    "document.json",
	}
}

// searchDirectoryForFiles searches a directory and its subdirectories for files with a specific name
func searchDirectoryForFiles(directory, filename string) ([]*yos.FilePathInfo, error) {
	dirs, err := yos.ListDir(directory)

	if err != nil {
		return nil, fmt.Errorf("got error listing directory: %w", err)
	}

	var result []*yos.FilePathInfo

	for _, dir := range dirs {
		files, err := yos.ListFile(dir.Path)

		if err != nil {
			return nil, fmt.Errorf("got error listing files: %w", err)
		}

		for _, f := range files {
			if f.Info.Name() == filename {
				result = append(result, f)
			}
		}
	}

	return result, nil
}

// loadBackupFile takes a slice of files and an interface, loads the json from disk and returns a slice of interfaces
func loadBackupFile(file *yos.FilePathInfo, i interface{}) error {
	dat, err := ioutil.ReadFile(file.Path)

	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Printf("Loaded backup file %s\n", file.Path)

	return json.Unmarshal(dat, i)
}
