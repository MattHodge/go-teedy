package teedy

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

type FileService struct {
	client   *resty.Client
	apiError *TeedyAPIError
}

func NewFileService(client *resty.Client, api string) *FileService {
	return &FileService{
		client:   client,
		apiError: NewTeedyAPIError(api),
	}
}

type File struct {
	Id         string `json:"id"`
	Processing bool   `json:"processing"`
	Name       string `json:"name"`
	Version    int    `json:"version"`
	MimeType   string `json:"mimetype"`
	DocumentId string `json:"document_id"`
	CreateDate *Time  `json:"create_date,omitempty"`
	Size       int    `json:"size"`
}

type FileList struct {
	Files []*File `json:"files,omitempty"`
}

type FileAddStatus struct {
	Status string `json:"status"`
	Id     string `json:"id"`
	Size   int    `json:"size"`
}

type ZippedFile struct {
	Filename string
	Content  []byte
}

func (f *FileService) GetAll() ([]*File, error) {
	resp, err := f.client.R().
		SetResult(&FileList{}).
		Get("api/file/list")

	err = checkRequestError(resp, err, f.apiError.GetAll)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*FileList).Files, nil
}

func (f *FileService) Get(documentId string) ([]*File, error) {
	resp, err := f.client.R().
		SetResult(&FileList{}).
		SetQueryParams(map[string]string{
			"id": documentId,
		}).
		Get("api/file/list")

	err = checkRequestError(resp, err, f.apiError.Get)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*FileList).Files, nil
}

func (f *FileService) GetData(id string) ([]byte, error) {
	resp, err := f.client.R().
		Get(fmt.Sprintf("api/file/%s/data", id))

	err = checkRequestError(resp, err, f.apiError.Custom("get file data"))

	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (f *FileService) Add(id, previousFileId string, file *os.File) (*FileAddStatus, error) {
	params := make(map[string]string)

	if len(id) > 0 {
		params["id"] = id
	}

	if len(previousFileId) > 0 {
		params["previousFileId"] = previousFileId
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("unable to get file information: %v", err)
	}

	resp, err := f.client.R().
		SetResult(&FileAddStatus{}).
		SetFormData(params).
		SetFileReader("file", fileInfo.Name(), file).
		Put("api/file")

	err = checkRequestError(resp, err, f.apiError.Add)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*FileAddStatus), nil
}

func (f *FileService) GetZippedFiles(id string) ([]byte, error) {
	resp, err := f.client.R().
		Get(fmt.Sprintf("api/file/zip?id=%s", id))

	err = checkRequestError(resp, err, f.apiError.Custom("getting zipped files"))

	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}
