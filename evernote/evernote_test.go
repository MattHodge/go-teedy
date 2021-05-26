package evernote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	ev, err := load("testdata/test.enex")
	require.NoError(t, err)

	assert.Len(t, ev.Notes, 2)
	assert.Equal(t, "1. test", ev.Notes[0].Title)
	assert.Equal(t, "2/2/2 test 2", ev.Notes[1].Title)
}
