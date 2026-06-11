package models

import (
	"github.com/vimek-go/pr-analizer-action/enum"
)

type Stats struct {
	Code     int
	Comments int
	Blanks   int
	Test     int
}

type LineType int

const (
	LineCode LineType = iota
	LineComment
	LineBlank
)

func (s *Stats) ToModelDiff(df enum.DiffType) DiffStats {
	if df == enum.DiffTypes.Removed() {
		return DiffStats{
			Code:            -s.Code,
			Comments:        -s.Comments,
			Blanks:          -s.Blanks,
			Test:            -s.Test,
			CodeRemoved:     s.Code,
			CommentsRemoved: s.Comments,
			BlanksRemoved:   s.Blanks,
			TestRemoved:     s.Test,
		}
	}
	return DiffStats{
		Code:          s.Code,
		Comments:      s.Comments,
		Blanks:        s.Blanks,
		Test:          s.Test,
		CodeAdded:     s.Code,
		CommentsAdded: s.Comments,
		BlanksAdded:   s.Blanks,
		TestAdded:     s.Test,
	}
}
