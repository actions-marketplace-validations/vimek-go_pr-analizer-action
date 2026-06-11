package models_test

import (
	"testing"

	"github.com/vimek-go/pr-analizer-action/enum"
	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/stretchr/testify/assert"
)

func TestStats_ToModelDiff_Added(t *testing.T) {
	t.Parallel()
	s := models.Stats{Code: 10, Comments: 3, Blanks: 2, Test: 1}
	diff := s.ToModelDiff(enum.DiffTypes.Added())

	assert.Equal(t, 10, diff.Code)
	assert.Equal(t, 10, diff.CodeAdded)
	assert.Equal(t, 0, diff.CodeRemoved)
	assert.Equal(t, 3, diff.CommentsAdded)
	assert.Equal(t, 0, diff.CommentsRemoved)
	assert.Equal(t, 2, diff.BlanksAdded)
	assert.Equal(t, 0, diff.BlanksRemoved)
	assert.Equal(t, 1, diff.TestAdded)
	assert.Equal(t, 0, diff.TestRemoved)
}

func TestStats_ToModelDiff_Removed(t *testing.T) {
	t.Parallel()
	s := models.Stats{Code: 10, Comments: 3, Blanks: 2, Test: 1}
	diff := s.ToModelDiff(enum.DiffTypes.Removed())

	assert.Equal(t, -10, diff.Code)
	assert.Equal(t, 0, diff.CodeAdded)
	assert.Equal(t, 10, diff.CodeRemoved)
	assert.Equal(t, 0, diff.CommentsAdded)
	assert.Equal(t, 3, diff.CommentsRemoved)
	assert.Equal(t, 0, diff.BlanksAdded)
	assert.Equal(t, 2, diff.BlanksRemoved)
	assert.Equal(t, 0, diff.TestAdded)
	assert.Equal(t, 1, diff.TestRemoved)
}
