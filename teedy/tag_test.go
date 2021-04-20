package teedy

import (
	"bytes"
	"github.com/MattHodge/go-teedy/mocks"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagService_GetAll(t *testing.T) {
	httpClient := &mocks.MockHTTPClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		const body = `
{
  "tags": [
    {
      "id": "d12a2911-c70f-4498-9c2b-19ba2646ca84",
      "name": "tax",
      "color": "#3a87ad"
    },
    {
      "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
      "name": "tax_2020",
      "color": "#3a87ad",
      "parent": "d12a2911-c70f-4498-9c2b-19ba2646ca84"
    }
  ]
}
`
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(body))),
		}, nil
	}}
	client := NewFakeClient(httpClient)

	tags, err := client.Tag.GetAll()
	assert.NoError(t, err)
	assert.Len(t, tags.Tags, 2)
}

func TestTagService_NewTag_Invalid_Color(t *testing.T) {
	tag, err := NewTag("test", "wrong", "")

	assert.Nil(t, tag)
	assert.Error(t, err, "a non-hex color for a tag should error")
}
func TestTagService_NewTag(t *testing.T) {
	tag, err := NewTag("test", "#ff0000", "")

	assert.NotNil(t, tag, "creating a new valid tag should not be nil")
	assert.NoError(t, err, "creating a valid tag should not cause an error")
}
