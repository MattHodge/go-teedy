package teedy_test

import (
	"testing"
	"time"

	"github.com/MattHodge/go-teedy/teedytest"

	"github.com/jarcoal/httpmock"

	"github.com/MattHodge/go-teedy/teedy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentService_GetAll(t *testing.T) {
	fixture := `
{
  "total": 2,
  "documents": [
    {
      "id": "557aeee4-aa21-4369-8c1c-9705c842c673",
      "highlight": null,
      "file_id": "3e9ea4c9-e34b-4ec7-a2dc-083bd669f52f",
      "title": "Insurance",
      "description": "<p><ul><li>foo</li><li>bar</li></ul></p>",
      "create_date": 1615849200000,
      "update_date": 1618911359554,
      "language": "eng",
      "shared": false,
      "active_route": false,
      "current_step_name": null,
      "file_count": 5,
      "tags": [
        {
          "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
          "name": "baz",
          "color": "#3a87ad"
        }
      ]
    },
    {
      "id": "c819adab-a297-4a83-8816-a19ce6d82d20",
      "highlight": null,
      "file_id": null,
      "title": "Blinds",
      "description": "<ul><li>fii</li><li>bii</li><li>bee</li></ul>",
      "create_date": 1615417200000,
      "update_date": 1618352722495,
      "language": "eng",
      "shared": false,
      "active_route": false,
      "current_step_name": null,
      "file_count": 0,
      "tags": []
    }
  ],
  "suggestions": []
}
`
	responder := teedytest.NewJsonResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://fake/api/document/list", responder)
	client := teedy.NewFakeClient()

	docs, err := client.Document.GetAll()
	require.NoError(t, err, "getting documents should not error")
	require.Len(t, docs, 2)
	expected := teedy.Timestamp{Time: time.Date(2021, 3, 15, 23, 00, 00, 00, time.UTC)}
	assert.Equal(t, expected.String(), docs[0].CreateDate.String(), "timestamp does not match expected")
}

func TestDocumentService_GetByTag(t *testing.T) {
	fixture := `
{
  "total": 2,
  "documents": [
    {
      "id": "1",
      "title": "Insurance",
      "tags": [
        {
          "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
          "name": "baz",
          "color": "#3a87ad"
        }
      ]
    },
    {
      "id": "2",
      "title": "Baz",
      "tags": [
        {
          "id": "bcd8e09b-84bc-4926-afce-222b7c21d8eb",
          "name": "baz",
          "color": "#3a87ad"
        }
      ]
    },
    {
      "id": "3",
      "title": "Blinds"
    }
  ],
  "suggestions": []
}
`
	responder := teedytest.NewJsonResponder(200, fixture)
	httpmock.RegisterResponder("GET", "http://fake/api/document/list", responder)
	client := teedy.NewFakeClient()

	docs, err := client.Document.GetByTagId("bcd8e09b-84bc-4926-afce-222b7c21d8eb")
	require.NoError(t, err, "getting documents by tag should not error")
	require.Len(t, docs, 2)
}

func TestDocumentService_AddDocument_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(teedytest.SkippingIntegrationMessage)
	}

	client := teedytest.SetupClient(t)

	doc := teedy.NewDocument("test document", "eng")

	createdTag, err := client.Document.Add(doc)
	require.NoError(t, err, "adding a new document should not cause error")
	assert.NotNil(t, createdTag.Id, "added document should be returned with an id")
}

func TestDocumentService_AddDocumentWithTag_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(teedytest.SkippingIntegrationMessage)
	}

	client := teedytest.SetupClient(t)

	tag1, err := client.Tag.Add(&teedy.Tag{
		Name:  "foo",
		Color: "#000000",
	})

	require.NoErrorf(t, err, "adding a tag should not fail")

	tag2, err := client.Tag.Add(&teedy.Tag{
		Name:  "bar",
		Color: "#ffffff",
	})

	require.NoErrorf(t, err, "adding a tag should not fail")

	doc := teedy.NewDocument("doc with tag 3", "eng")
	doc.Tags = []*teedy.Tag{
		tag1,
		tag2,
	}

	require.NotNilf(t, doc, "creating a new valid document should not be nil")

	createdDoc, err := client.Document.Add(doc)
	require.NoErrorf(t, err, "adding a new document should not cause error")
	assert.NotNil(t, createdDoc.Id, "added document should be returned with an id")

	readCreatedDoc, err := client.Document.Get(createdDoc.Id)
	require.NoErrorf(t, err, "getting the created document should not cause error")

	assert.Len(t, doc.Tags, len(readCreatedDoc.Tags), "tags on new document should match created doc")
}

func TestDocumentService_UpdateTagIDsByName(t *testing.T) {
	document := teedy.NewDocument("foo", "eng")

	oldTags := []*teedy.Tag{
		{
			Id:   "aaa",
			Name: "foo",
		},
		{
			Id:   "bbb",
			Name: "bar",
		},
		{
			Id:   "ccc",
			Name: "baz",
		},
	}

	document.Tags = oldTags

	newTags := []*teedy.Tag{
		{
			Id:   "100",
			Name: "bar",
		},
		{
			Id:   "200",
			Name: "foo",
		},
	}

	document.UpdateTagIDsByName(newTags)

	want := []*teedy.Tag{
		{
			Id:   "200",
			Name: "foo",
		},
		{
			Id:   "100",
			Name: "bar",
		},
		{
			Id:   "ccc",
			Name: "baz",
		},
	}

	assert.EqualValues(t, want, document.Tags, "documents tag IDs should have been updated with passed values")
}
