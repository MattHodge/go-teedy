package backup

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/MattHodge/go-teedy/teedy"
)

type FileBackup struct {
	FullDirectory string
	File          *teedy.File
	client        *teedy.Client
}

func (f *FileBackup) Save() error {
	os.MkdirAll(f.FullDirectory, 0700)

	fbytes, err := f.client.File.GetData(f.File.Id)

	if err != nil {
		return fmt.Errorf("can't get data file with id %s", f.File.Id)
	}

	fullFilePath := filepath.Join(f.FullDirectory, f.File.Name)

	err = ioutil.WriteFile(fullFilePath, fbytes, 0644)

	if err != nil {
		return fmt.Errorf("can't write file: %w", err)
	}

	fmt.Printf("File exported succesfully: %s\n", fullFilePath)

	return nil
}

func File(client *teedy.Client, file *teedy.File, basePath string) *FileBackup {
	fullDirectory := filepath.Join(basePath, "documents", file.DocumentId, "files")
	return &FileBackup{
		FullDirectory: fullDirectory,
		File:          file,
		client:        client,
	}
}
