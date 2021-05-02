package teedy_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/MattHodge/go-teedy/mocks"
	"github.com/MattHodge/go-teedy/teedy"
	"github.com/stretchr/testify/require"

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
	client := teedy.NewFakeClient(httpClient)

	tags, err := client.Tag.GetAll()
	require.NoError(t, err)
	assert.Len(t, tags.Tags, 2)
}

func TestTagService_NewTag_Invalid_Color(t *testing.T) {
	tag, err := teedy.NewTag("test", "wrong", "")

	assert.Nil(t, tag)
	assert.Error(t, err, "a non-hex color for a tag should error")
}

func TestTagService_NewTag(t *testing.T) {
	tag, err := teedy.NewTag("test", "#ff0000", "")

	require.NotNil(t, tag, "creating a new valid tag should not be nil")
	assert.NoError(t, err, "creating a valid tag should not cause an error")
}

func TestTagService_AddTag_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)

	tag, err := teedy.NewTag("test", "#fff000", "")
	require.NotNil(t, tag, "creating a new valid tag should not be nil")

	createdTag, err := client.Tag.Add(tag)
	require.NoError(t, err, "adding a new tag should not cause error")
	assert.NotNil(t, createdTag.Id, "added tag should be returned with an id")
}

func TestTagService_DeleteTag_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	client := setup(t)

	tag, err := teedy.NewTag("test", "#fff000", "")
	require.NotNil(t, tag, "creating a new valid tag should not be nil")

	createdTag, err := client.Tag.Add(tag)
	require.NoError(t, err, "adding a new tag should not cause error")

	tagDeleteStatus, err := client.Tag.Delete(createdTag.Id)

	require.NoError(t, err, "deleting tag should not error")
	assert.Equal(t, tagDeleteStatus.Status, "ok", "tag delete status should be ok")
}

func TestTagService_UpdateTag_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(testSkippingIntegrationTest)
	}

	const (
		updatedTagColor = "#ff7816"
		updatedTagName  = "biz"
	)

	client := setup(t)

	tag, err := teedy.NewTag("test", "#fff000", "")
	require.NotNil(t, tag, "creating a new valid tag should not be nil")

	createdTag, err := client.Tag.Add(tag)
	require.NoError(t, err, "adding a new tag should not cause error")

	tagNewColor, err := teedy.NewTag(updatedTagName, updatedTagColor, "")

	updatedTag, err := client.Tag.Update(createdTag.Id, tagNewColor)

	require.NoError(t, err, "updating tag should not error")

	updatedTagReloadedFromAPI, err := client.Tag.Get(updatedTag.Id)
	assert.NoError(t, err, "getting updated tag should not return error")
	assert.Equal(t, updatedTagReloadedFromAPI.Color, updatedTagColor)
	assert.Equal(t, updatedTagReloadedFromAPI.Name, updatedTagName)
}
