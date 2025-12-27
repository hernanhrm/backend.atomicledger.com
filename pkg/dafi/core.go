// Package dafi provides dynamic filtering, sorting, and pagination capabilities.
package dafi

// Criteria defines the complete set of query criteria including selection, joins, filters, sorting, and pagination.
type Criteria struct {
	SelectColumns   []string
	Joins           []string
	Filters         Filters
	FiltersByModule map[string]Filters
	Sorts           Sorts
	Pagination      Pagination
}

// New creates a new empty Criteria.
func New() Criteria {
	return Criteria{}
}

// Where creates a new Criteria with an initial filter.
func Where(field string, operator FilterOperator, value any) Criteria {
	return Criteria{
		Filters: FilterBy(field, operator, value),
	}
}

// AndGroup adds a group of filters combined with AND logic to the current filters.
func (c Criteria) AndGroup(filters ...Filter) Criteria {
	c.Filters = c.Filters.AndGroup(filters...)

	return c
}

// OrGroup adds a group of filters combined with OR logic to the current filters.
func (c Criteria) OrGroup(filters ...Filter) Criteria {
	c.Filters = c.Filters.OrGroup(filters...)

	return c
}

// Or adds a new filter combined with OR logic.
func (c Criteria) Or(field string, operator FilterOperator, value any) Criteria {
	c.Filters = c.Filters.Or(field, operator, value)

	return c
}

// And adds a new filter combined with AND logic.
func (c Criteria) And(field string, operator FilterOperator, value any) Criteria {
	c.Filters = c.Filters.And(field, operator, value)

	return c
}

// SortBy adds a sort instruction.
func (c Criteria) SortBy(field string, sortType SortType) Criteria {
	c.Sorts = append(c.Sorts, Sort{Field: SortBy(field), Type: sortType})

	return c
}

// Limit sets the page size limit.
func (c Criteria) Limit(value uint) Criteria {
	c.Pagination.PageSize = value

	return c
}

// Page sets the page number.
func (c Criteria) Page(value uint) Criteria {
	c.Pagination.PageNumber = value

	return c
}

// Select sets the columns to be selected.
func (c Criteria) Select(columns ...string) Criteria {
	c.SelectColumns = columns

	return c
}
