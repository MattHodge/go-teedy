package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/MattHodge/go-teedy/teedy"
)

func (c *Client) File(file *teedy.File) error {
	fileDirectory := filepath.Join(c.DocumentBackupDirectory(file.DocumentId), "files")
	os.MkdirAll(fileDirectory, 0700)

	fbytes, err := c.client.File.GetData(file.Id)

	if err != nil {
		return fmt.Errorf("can't get data file with id %s", file.Id)
	}

	fullFilePath := filepath.Join(fileDirectory, file.Name)

	err = ioutil.WriteFile(fullFilePath, fbytes, 0644)

	if err != nil {
		return fmt.Errorf("can't write file: %w", err)
	}

	return nil
}
