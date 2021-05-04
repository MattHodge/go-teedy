package teedy

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type DocumentService struct {
	client   *resty.Client
	apiError *TeedyAPIError
}

func NewDocumentService(client *resty.Client, api string) *DocumentService {
	return &DocumentService{
		client:   client,
		apiError: NewTeedyAPIError(api),
	}
}

type DocumentList struct {
	Documents   []Document `json:"documents"`
	Total       int        `json:"total"`
	Suggestions []string   `json:"suggestions"`
}

type ACLs struct {
	Id   string `json:"id"`
	Perm string `json:"perm"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type InheritedACLs struct {
	Name        string `json:"name"`
	Perm        string `json:"perm"`
	SourceId    string `json:"source_id"`
	SourceName  string `json:"source_name"`
	SourceColor string `json:"source_color"`
	Type        string `json:"type"`
	Id          string `json:"id"`
}

type Contributors struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Relations struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Source string `json:"source"`
}

type RouteStep struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Transitionable bool   `json:"transitionable"`
}

type Document struct {
	Id              string          `json:"id,omitempty"`
	Highlight       string          `json:"highlight,omitempty"`
	FileID          string          `json:"file_id,omitempty"`
	Title           string          `json:"title"`
	Description     string          `json:"description,omitempty"`
	CreateDate      Time            `json:"create_date,omitempty"`
	UpdateDate      Time            `json:"update_date,omitempty"`
	Language        string          `json:"language"`
	Shared          bool            `json:"shared,omitempty"`
	ActiveRoute     bool            `json:"active_route,omitempty"`
	CurrentStepName bool            `json:"current_step_name,omitempty"`
	FileCount       int             `json:"file_count,omitempty"`
	Tags            []Tag           `json:"tags,omitempty"`
	Subject         string          `json:"subject,omitempty"`
	Identifier      string          `json:"identifier,omitempty"`
	Publisher       string          `json:"publisher,omitempty"`
	Format          string          `json:"format,omitempty"`
	Source          string          `json:"source,omitempty"`
	Type            string          `json:"type,omitempty"`
	Coverage        string          `json:"coverage,omitempty"`
	Rights          string          `json:"rights,omitempty"`
	Creator         string          `json:"creator,omitempty"`
	Writeable       bool            `json:"writeable,omitempty"`
	ACLs            []ACLs          `json:"acls,omitempty"`
	InheritedACLs   []InheritedACLs `json:"inherited_acls,omitempty"`
	Contributors    []Contributors  `json:"contributors,omitempty"`
	Relations       []Relations     `json:"relations,omitempty"`
	RouteStep       RouteStep       `json:"route_step,omitempty"`
}

type DocumentDeleteStatus struct {
	Status string `json:"Status"`
}

func NewDocument(title, language string) (*Document, error) {
	return &Document{
		Title:    title,
		Language: language,
	}, nil
}

func (t *DocumentService) GetAll() (*DocumentList, error) {
	resp, err := t.client.R().
		SetResult(&DocumentList{}).
		Get("api/document/list")

	err = checkRequestError(resp, err, t.apiError.GetAll)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*DocumentList), nil
}

func (t *DocumentService) Get(id string) (*Document, error) {
	resp, err := t.client.R().
		SetResult(&Document{}).
		Get(fmt.Sprintf("api/document/%s", id))

	err = checkRequestError(resp, err, t.apiError.Get)

	return resp.Result().(*Document), nil
}

func (t *DocumentService) Add(d *Document) (*Document, error) {
	// builds the form data for creating a document in the teedy api
	fv := NewFormValues()
	fv.AddMandatory("title", d.Title)
	fv.AddMandatory("language", d.Language)
	fv.AddIfNotEmpty("title", d.Title)
	fv.AddIfNotEmpty("description", d.Description)
	fv.AddIfNotEmpty("subject", d.Subject)
	fv.AddIfNotEmpty("identifier", d.Identifier)
	fv.AddIfNotEmpty("publisher", d.Publisher)
	fv.AddIfNotEmpty("format", d.Format)
	fv.AddIfNotEmpty("source", d.Source)
	fv.AddIfNotEmpty("type", d.Type)
	fv.AddIfNotEmpty("coverage", d.Coverage)
	fv.AddIfNotEmpty("rights", d.Rights)

	for _, tag := range d.Tags {
		fv.AddIfNotEmpty("tag", tag.Id)
	}

	for _, rel := range d.Relations {
		fv.AddIfNotEmpty("relations", rel.Id)
	}

	body, err := fv.Result()

	resp, err := t.client.R().
		SetResult(&Document{}).
		SetFormDataFromValues(body).
		Put("api/document")

	err = checkRequestError(resp, err, t.apiError.Add)

	return resp.Result().(*Document), nil
}

func (t *DocumentService) Delete(id string) (*DocumentDeleteStatus, error) {
	resp, err := t.client.R().
		SetResult(&DocumentDeleteStatus{}).
		Delete(fmt.Sprintf("api/document/%s", id))

	err = checkRequestError(resp, err, t.apiError.Delete)

	return resp.Result().(*DocumentDeleteStatus), nil
}
