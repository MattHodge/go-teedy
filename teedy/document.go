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
	Documents   []*Document `json:"documents"`
	Total       int         `json:"total"`
	Suggestions []string    `json:"suggestions"`
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
	Name           string `json:"name,omitempty"`
	Type           string `json:"type,omitempty"`
	Transitionable bool   `json:"transitionable,omitempty"`
}

type Document struct {
	Id              string           `json:"id,omitempty"`
	Highlight       string           `json:"highlight,omitempty"`
	FileID          string           `json:"file_id,omitempty"`
	Title           string           `json:"title"`
	Description     string           `json:"description,omitempty"`
	CreateDate      *Timestamp       `json:"create_date,omitempty"`
	UpdateDate      *Timestamp       `json:"update_date,omitempty"`
	Language        string           `json:"language"`
	Shared          bool             `json:"shared,omitempty"`
	ActiveRoute     bool             `json:"active_route,omitempty"`
	CurrentStepName bool             `json:"current_step_name,omitempty"`
	FileCount       int              `json:"file_count,omitempty"`
	Tags            []*Tag           `json:"tags,omitempty"`
	Subject         string           `json:"subject,omitempty"`
	Identifier      string           `json:"identifier,omitempty"`
	Publisher       string           `json:"publisher,omitempty"`
	Format          string           `json:"format,omitempty"`
	Source          string           `json:"source,omitempty"`
	Type            string           `json:"type,omitempty"`
	Coverage        string           `json:"coverage,omitempty"`
	Rights          string           `json:"rights,omitempty"`
	Creator         string           `json:"creator,omitempty"`
	Writeable       bool             `json:"writeable,omitempty"`
	ACLs            []*ACLs          `json:"acls,omitempty"`
	InheritedACLs   []*InheritedACLs `json:"inherited_acls,omitempty"`
	Contributors    []*Contributors  `json:"contributors,omitempty"`
	Relations       []*Relations     `json:"relations,omitempty"`
	RouteStep       *RouteStep       `json:"route_step,omitempty"`
}

// UpdateTagIDsByName updates the documents tag ID's if the incoming tag matches its name. This is useful
// when importing backups of documents
func (d *Document) UpdateTagIDsByName(tags []*Tag) {
	for _, documentTag := range d.Tags {
		for _, incomingTag := range tags {
			if documentTag.Name == incomingTag.Name {
				documentTag.Id = incomingTag.Id
			}
		}
	}
}

type DocumentDeleteStatus struct {
	Status string `json:"Status"`
}

func NewDocument(title, language string) *Document {
	return &Document{
		Title:    title,
		Language: language,
	}
}

func (d *DocumentService) GetAll() (*DocumentList, error) {
	resp, err := d.client.R().
		SetResult(&DocumentList{}).
		Get("api/document/list")

	err = checkRequestError(resp, err, d.apiError.GetAll)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*DocumentList), nil
}

func (d *DocumentService) Get(id string) (*Document, error) {
	resp, err := d.client.R().
		SetResult(&Document{}).
		Get(fmt.Sprintf("api/document/%s", id))

	err = checkRequestError(resp, err, d.apiError.Get)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Document), nil
}

func (d *DocumentService) GetByTitle(title string) (*Document, error) {
	docs, err := d.GetAll()

	if err != nil {
		return nil, err
	}

	for _, doc := range docs.Documents {
		if doc.Title == title {
			return doc, nil
		}
	}

	return nil, nil
}

func (d *DocumentService) Add(doc *Document) (*Document, error) {
	// builds the form data for creating a document in the teedy api
	fv := NewFormValues()
	fv.AddMandatory("title", doc.Title)
	fv.AddMandatory("language", doc.Language)
	fv.AddIfNotEmpty("title", doc.Title)
	fv.AddIfNotEmpty("description", doc.Description)
	fv.AddIfNotEmpty("subject", doc.Subject)
	fv.AddIfNotEmpty("identifier", doc.Identifier)
	fv.AddIfNotEmpty("publisher", doc.Publisher)
	fv.AddIfNotEmpty("format", doc.Format)
	fv.AddIfNotEmpty("source", doc.Source)
	fv.AddIfNotEmpty("type", doc.Type)
	fv.AddIfNotEmpty("coverage", doc.Coverage)
	fv.AddIfNotEmpty("rights", doc.Rights)

	if doc.CreateDate != nil {
		fv.AddIfNotEmpty("create_date", doc.CreateDate.Marshal())
	}

	for _, tag := range doc.Tags {
		fv.AddIfNotEmpty("tags", tag.Id)
	}

	for _, rel := range doc.Relations {
		fv.AddIfNotEmpty("relations", rel.Id)
	}

	body, err := fv.Result()

	resp, err := d.client.R().
		SetResult(&Document{}).
		SetFormDataFromValues(body).
		Put("api/document")

	err = checkRequestError(resp, err, d.apiError.Add)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Document), nil
}

func (d *DocumentService) Delete(id string) (*DocumentDeleteStatus, error) {
	resp, err := d.client.R().
		SetResult(&DocumentDeleteStatus{}).
		Delete(fmt.Sprintf("api/document/%s", id))

	err = checkRequestError(resp, err, d.apiError.Delete)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*DocumentDeleteStatus), nil
}
