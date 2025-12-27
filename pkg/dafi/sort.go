package dafi

// SortType defines the sort direction.
type SortType string

const (
	// Asc represents ascending sort order.
	Asc SortType = "ASC"
	// Desc represents descending sort order.
	Desc SortType = "DESC"
	// None represents no sort order.
	None SortType = ""
)

// SortBy represents the field to sort by.
type SortBy string

// Sort represents a sort instruction.
type Sort struct {
	Field SortBy
	Type  SortType
}

// Sorts represents a collection of sort instructions.
type Sorts []Sort

// IsZero checks if the sorts collection is empty.
func (s Sorts) IsZero() bool {
	return len(s) == 0
}
