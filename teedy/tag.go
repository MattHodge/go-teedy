package teedy

import (
	"fmt"
	"net/url"
	"regexp"
)

type TagService struct {
	client *Client
}

func NewTagService(client *Client) *TagService {
	return &TagService{client: client}
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

func (t *TagService) GetAll() (*TagList, error) {
	endpoint := "api/tag/list"
	tags, err := t.client.request(endpoint, "GET", nil, new(TagList))

	if err != nil {
		return nil, fmt.Errorf("error getting all tags: %v", err)
	}

	return tags.(*TagList), nil
}

func (t *TagService) Get(id string) (*Tag, error) {
	endpoint := fmt.Sprintf("api/tag/%s", id)
	tag, err := t.client.request(endpoint, "GET", nil, new(Tag))

	if err != nil {
		return nil, fmt.Errorf("error getting tag: %v", err)
	}

	return tag.(*Tag), nil
}

func (t *TagService) Add(tag *Tag) (*Tag, error) {
	endpoint := "api/tag"

	formData := url.Values{
		"name":   {tag.Name},
		"color":  {tag.Color},
		"parent": {tag.Parent},
	}

	returnTag, err := t.client.request(endpoint, "PUT", formData, new(Tag))

	if err != nil {
		return nil, fmt.Errorf("error adding tag: %v", err)
	}

	return returnTag.(*Tag), nil
}
