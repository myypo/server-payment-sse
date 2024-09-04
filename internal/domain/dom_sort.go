package dom

import "fmt"

type SortByTime int

const (
	CreatedAt SortByTime = iota
	UpdatedAt
)

var validSortByTimeStrings [2]string = [2]string{
	"created_at",
	"updated_at",
}

func (s SortByTime) String() string {
	return validSortByTimeStrings[s]
}

func SortByTimeFromString(s string) (SortByTime, error) {
	switch s {
	case CreatedAt.String():
		return CreatedAt, nil
	case UpdatedAt.String():
		return UpdatedAt, nil
	default:
		return -1, fmt.Errorf(
			"unknown `sort by time` provided: %s, valid options are: %s",
			s,
			validSortByTimeStrings,
		)
	}
}

type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

var validSortOrderStrings [2]string = [2]string{
	"asc",
	"desc",
}

func (s SortOrder) String() string {
	return validSortOrderStrings[s]
}

func SortOrderFromString(s string) (SortOrder, error) {
	switch s {
	case Ascending.String():
		return Ascending, nil
	case Descending.String():
		return Descending, nil
	default:
		return -1, fmt.Errorf(
			"unknown `sort order` provided: %s, valid options are: %s",
			s,
			validSortOrderStrings,
		)
	}
}
