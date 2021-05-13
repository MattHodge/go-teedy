package restore_test

import (
	"path/filepath"
	"testing"

	"github.com/MattHodge/go-teedy/restore"

	"github.com/MattHodge/go-teedy/teedy"
	"github.com/MattHodge/go-teedy/teedytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	tc := teedy.NewFakeClient()
	backupDir := t.TempDir()
	restoreClient := restore.NewRestoreClient(tc, backupDir)

	tag1 := `
{
    "id": "aa",
    "name": "tax_2021",
    "color": "#00b046"
}
`

	tag2 := `
{
    "id": "bb",
    "name": "tax_2021",
    "color": "#00b046"
}
`

	teedytest.WriteToFile(t, filepath.Join(backupDir, "tags", "aa", "tag.json"), tag1)
	teedytest.WriteToFile(t, filepath.Join(backupDir, "tags", "bb", "tag.json"), tag2)

	tags, err := restoreClient.ViewTags()

	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "aa", tags[0].Id)
	assert.Equal(t, "bb", tags[1].Id)
}

func TestTags_UpdateParentTag_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip(teedytest.SkippingIntegrationMessage)
	}

	client := teedytest.SetupClient(t)

	defer cleanTag(t, client, "parent")
	defer cleanTag(t, client, "child")
	defer cleanTag(t, client, "nested-child")

	backupDir := t.TempDir()
	restoreClient := restore.NewRestoreClient(client, backupDir)

	tag1 := `
{
    "id": "aa",
    "name": "parent",
    "color": "#000000"
}
`

	tag2 := `
{
    "id": "bb",
    "name": "child",
    "color": "#000000",
	"parent": "aa" 
}
`

	tag3 := `
{
    "id": "cc",
    "name": "nested-child",
    "color": "#000000",
	"parent": "bb" 
}
`

	teedytest.WriteToFile(t, filepath.Join(backupDir, "tags", "aa", "tag.json"), tag1)
	teedytest.WriteToFile(t, filepath.Join(backupDir, "tags", "bb", "tag.json"), tag2)
	teedytest.WriteToFile(t, filepath.Join(backupDir, "tags", "cc", "tag.json"), tag3)

	err := restoreClient.Tags()

	require.NoError(t, err)

	parentTag := getTagByName(t, client, "parent")
	childTag := getTagByName(t, client, "child")
	nestedChildTag := getTagByName(t, client, "nested-child")

	assert.Equal(t, parentTag.Id, childTag.Parent, "the child tag should reference the parent tags id")
	assert.Equal(t, childTag.Id, nestedChildTag.Parent, "the nested-child tag should reference the child tags id")
}

func cleanTag(t *testing.T, client *teedy.Client, tagName string) {
	tag, err := client.Tag.GetByName(tagName)

	require.NoError(t, err, "shouldn't return error getting tag")

	if tag != nil && tag.Id != "" {
		_, err := client.Tag.Delete(tag.Id)

		require.NoError(t, err, "shouldn't return error deleting tag")
	}
}

func getTagByName(t *testing.T, client *teedy.Client, tagName string) *teedy.Tag {
	tag, err := client.Tag.GetByName(tagName)

	require.NoError(t, err, "shouldn't return error getting tag")

	return tag
}
