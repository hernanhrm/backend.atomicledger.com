package dafi

// Pagination defines the pagination parameters.
type Pagination struct {
	PageNumber uint
	PageSize   uint
}

// IsZero checks if the pagination is empty.
func (p Pagination) IsZero() bool {
	return p.PageNumber == 0 && p.PageSize == 0
}

// HasPageNumber checks if the page number is set.
func (p Pagination) HasPageNumber() bool {
	return p.PageNumber > 0
}

// HasPageSize checks if the page size is set.
func (p Pagination) HasPageSize() bool {
	return p.PageSize > 0
}
