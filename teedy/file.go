package teedy

import (
	"fmt"
	"os"
)

type FileService struct {
	client *Client
}

func NewFileService(client *Client) *FileService {
	return &FileService{client: client}
}

type File struct {
	Id         string `json:"id"`
	Processing bool   `json:"processing"`
	Name       string `json:"name"`
	Version    int    `json:"version"`
	MimeType   string `json:"mimetype"`
	DocumentId string `json:"document_id"`
	CreateDate Time   `json:"create_date,omitempty"`
	Size       int    `json:"size"`
}

type FileList struct {
	Files []File `json:"files,omitempty"`
}

type FileAddStatus struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Size   int    `json:"size"`
}

func (f *FileService) GetAll() (*FileList, error) {
	endpoint := "api/file/list"
	files, err := f.client.requestUnmarshal(endpoint, "GET", nil, new(FileList))

	if err != nil {
		return nil, fmt.Errorf("error getting all files: %v", err)
	}

	return files.(*FileList), nil
}

func (f *FileService) Add(id, previousFileId string, file *os.File) (*FileAddStatus, error) {
	endpoint := "api/file"

	params := make(map[string]string)

	if len(id) > 0 {
		params["id"] = id
	}

	if len(previousFileId) > 0 {
		params["previousFileId"] = previousFileId
	}

	resp, err := f.client.multipartUpload(endpoint, "PUT", params, file)

	if err != nil {
		return nil, fmt.Errorf("error adding file: %v", err)
	}

	fileStatus, err := unmarshalResponse(resp, new(FileAddStatus))

	if err != nil {
		return nil, fmt.Errorf("error unmarshelling file add status: %v", err)
	}

	return fileStatus.(*FileAddStatus), nil
}

func (f *FileService) GetZippedFiles(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("api/file/zip?id=%s", id)

	bytes, err := f.client.request(endpoint, "GET", nil)

	if err != nil {
		return nil, fmt.Errorf("error getting zipped file: %v", err)
	}

	return bytes, nil
}
