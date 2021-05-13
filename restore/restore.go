package restore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/1set/gut/yos"
	"github.com/MattHodge/go-teedy/teedy"
)

// SearchDirectoryForBackupFiles searches a directory and its subdirectories for files with a specific name
func SearchDirectoryForBackupFiles(directory, filename string) ([]*yos.FilePathInfo, error) {
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

// LoadBackupFile takes a slice of files and an interface, loads the json from disk and returns a slice of interfaces
func LoadBackupFile(file *yos.FilePathInfo, i interface{}) error {
	dat, err := ioutil.ReadFile(file.Path)

	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Printf("Loaded backup file %s\n", file.Path)

	return json.Unmarshal(dat, i)
}

func LoadBackupTags(files []*yos.FilePathInfo, filename string) ([]*teedy.Tag, error) {
	var res []*teedy.Tag

	for _, file := range files {
		if file.Info.Name() != filename {
			return nil, fmt.Errorf("won't load non-tag files: %s", file.Path)
		}
		tag := new(teedy.Tag)
		err := LoadBackupFile(file, tag)

		if err != nil {
			return nil, fmt.Errorf("unable to load backup file: %w", err)
		}

		res = append(res, tag)
	}

	return res, nil
}

func LoadBackupDocuments(directory string) ([]*teedy.Document, error) {
	files, err := SearchDirectoryForBackupFiles(directory, "document.json")

	if err != nil {
		return nil, err
	}

	var res []*teedy.Document

	for _, file := range files {
		doc := new(teedy.Document)
		err := LoadBackupFile(file, doc)

		if err != nil {
			return nil, fmt.Errorf("unable to load backup file: %w", err)
		}

		res = append(res, doc)
	}

	return res, nil
}

func Tags(client *teedy.Client, tags []*teedy.Tag) error {
	for _, t := range tags {
		// check if tag exists
		existingTag, err := client.Tag.GetByName(t.Name)

		if existingTag != nil {
			// delete it incase there are changes.. could be an update
			_, err := client.Tag.Delete(existingTag.Id)

			if err != nil {
				return err
			}
		}

		_, err = client.Tag.Add(t)

		if err != nil {
			return err
		}
	}

	return nil
}

func Documents(client *teedy.Client, docs []*teedy.Document) error {
	for _, d := range docs {
		// check if doc exists
		existingDoc, err := client.Document.GetByTitle(d.Title)

		if existingDoc != nil {
			// delete it incase there are changes.. could be an update
			_, err := client.Document.Delete(existingDoc.Id)

			if err != nil {
				return err
			}
		}

		_, err = client.Document.Add(d)

		if err != nil {
			return err
		}
	}

	return nil
}
