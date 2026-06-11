package enum

type DiffType string

const (
	diffTypeAdded   DiffType = "added"
	diffTypeRemoved DiffType = "removed"
)

func (d DiffType) String() string {
	return string(d)
}

type diffType struct{}

func (diffType) Added() DiffType {
	return diffTypeAdded
}

func (diffType) Removed() DiffType {
	return diffTypeRemoved
}

var DiffTypes = diffType{}
