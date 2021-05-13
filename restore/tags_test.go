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
