package models_test

import (
	"testing"

	"github.com/vimek-go/pr-analizer-action/models"

	"github.com/stretchr/testify/assert"
)

func TestDiffReport_Add_ZeroValue(t *testing.T) {
	t.Parallel()
	var r models.DiffReport

	assert.NotPanics(t, func() {
		r.Add("Go", models.DiffStats{Code: 5, CodeAdded: 5})
	})
	assert.Equal(t, 5, r.Total.Code)
	assert.Equal(t, 5, r.ByLanguage["Go"].Code)
}

func TestDiffReport_Add_AccumulatesTotal(t *testing.T) {
	t.Parallel()
	var r models.DiffReport

	r.Add("Go", models.DiffStats{Code: 10, CodeAdded: 12, CodeRemoved: 2})
	r.Add("Python", models.DiffStats{Code: 5, CodeAdded: 5})

	assert.Equal(t, 15, r.Total.Code)
	assert.Equal(t, 17, r.Total.CodeAdded)
	assert.Equal(t, 2, r.Total.CodeRemoved)
}

func TestDiffReport_Add_ByLanguage(t *testing.T) {
	t.Parallel()
	var r models.DiffReport

	r.Add("Go", models.DiffStats{CodeAdded: 5})
	r.Add("Go", models.DiffStats{CodeAdded: 3})

	assert.Equal(t, 8, r.ByLanguage["Go"].CodeAdded)
}

func TestDiffReport_Add_EmptyLangSkipsMap(t *testing.T) {
	t.Parallel()
	var r models.DiffReport

	r.Add("", models.DiffStats{Code: 5})

	assert.Equal(t, 5, r.Total.Code)
	assert.Empty(t, r.ByLanguage)
}
