package teedy

import (
	"fmt"
	"regexp"

	"github.com/go-resty/resty/v2"
)

type TagService struct {
	client   *resty.Client
	apiError *TeedyAPIError
}

func NewTagService(client *resty.Client, api string) *TagService {
	return &TagService{
		client:   client,
		apiError: NewTeedyAPIError(api),
	}
}

type TagList struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	Parent string `json:"parent,omitempty"`
}

type TagDeleteStatus struct {
	Status string `json:"status"`
}

func NewTag(name, color, parent string) (*Tag, error) {
	re := regexp.MustCompile(`^#([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$`)

	if !re.MatchString(color) {
		return nil, fmt.Errorf("color needs to be a hex code like #ff0000, got %s", color)
	}

	return &Tag{
		Name:   name,
		Color:  color,
		Parent: parent,
	}, nil
}

func (t *TagService) GetAll() ([]Tag, error) {
	resp, err := t.client.R().
		SetResult(&TagList{}).
		Get("api/tag/list")

	err = checkRequestError(resp, err, t.apiError.GetAll)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*TagList).Tags, nil
}

func (t *TagService) Get(id string) (*Tag, error) {
	resp, err := t.client.R().
		SetResult(&Tag{}).
		Get(fmt.Sprintf("api/tag/%s", id))

	err = checkRequestError(resp, err, t.apiError.Get)

	return resp.Result().(*Tag), nil
}

func (t *TagService) Add(tag *Tag) (*Tag, error) {
	resp, err := t.client.R().
		SetResult(&Tag{}).
		SetFormData(map[string]string{
			"name":   tag.Name,
			"color":  tag.Color,
			"parent": tag.Parent,
		}).
		Put("api/tag")

	err = checkRequestError(resp, err, t.apiError.Add)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*Tag), nil
}

func (t *TagService) Delete(id string) (*TagDeleteStatus, error) {
	resp, err := t.client.R().
		SetResult(&TagDeleteStatus{}).
		Delete(fmt.Sprintf("api/tag/%s", id))

	err = checkRequestError(resp, err, t.apiError.Delete)

	return resp.Result().(*TagDeleteStatus), nil
}

func (t *TagService) Update(id string, tag *Tag) (*Tag, error) {
	resp, err := t.client.R().
		SetResult(&Tag{}).
		SetFormData(map[string]string{
			"name":   tag.Name,
			"color":  tag.Color,
			"parent": tag.Parent,
		}).
		Post(fmt.Sprintf("api/tag/%s", id))

	err = checkRequestError(resp, err, t.apiError.Update)

	return resp.Result().(*Tag), nil
}
