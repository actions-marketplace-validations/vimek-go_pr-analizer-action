package models

type DiffStats struct {
	Code     int
	Comments int
	Blanks   int
	Test     int

	CodeAdded       int
	CodeRemoved     int
	CommentsAdded   int
	CommentsRemoved int
	BlanksAdded     int
	BlanksRemoved   int
	TestAdded       int
	TestRemoved     int
}

type DiffReport struct {
	Total      DiffStats
	ByLanguage map[string]DiffStats
}

func (r *DiffReport) Add(langName string, diff DiffStats) {
	r.Total.Code += diff.Code
	r.Total.Comments += diff.Comments
	r.Total.Blanks += diff.Blanks
	r.Total.Test += diff.Test

	r.Total.CodeAdded += diff.CodeAdded
	r.Total.CodeRemoved += diff.CodeRemoved
	r.Total.CommentsAdded += diff.CommentsAdded
	r.Total.CommentsRemoved += diff.CommentsRemoved
	r.Total.BlanksAdded += diff.BlanksAdded
	r.Total.BlanksRemoved += diff.BlanksRemoved
	r.Total.TestAdded += diff.TestAdded
	r.Total.TestRemoved += diff.TestRemoved

	if langName != "" {
		if r.ByLanguage == nil {
			r.ByLanguage = make(map[string]DiffStats)
		}
		langStats := r.ByLanguage[langName]
		langStats.Code += diff.Code
		langStats.Comments += diff.Comments
		langStats.Blanks += diff.Blanks
		langStats.Test += diff.Test

		langStats.CodeAdded += diff.CodeAdded
		langStats.CodeRemoved += diff.CodeRemoved
		langStats.CommentsAdded += diff.CommentsAdded
		langStats.CommentsRemoved += diff.CommentsRemoved
		langStats.BlanksAdded += diff.BlanksAdded
		langStats.BlanksRemoved += diff.BlanksRemoved
		langStats.TestAdded += diff.TestAdded
		langStats.TestRemoved += diff.TestRemoved

		r.ByLanguage[langName] = langStats
	}
}
