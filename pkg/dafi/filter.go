package dafi

type (
	// FilterField represents the name of the field to filter by.
	FilterField string
	// FilterValue represents the value to filter by.
	FilterValue any
)

// FilterOperator represents the comparison operator for a filter.
type FilterOperator string

const (
	// Equal matches values equal to the specified value.
	Equal FilterOperator = "eq"
	// NotEqual matches values not equal to the specified value.
	NotEqual FilterOperator = "ne"
	// Greater matches values greater than the specified value.
	Greater FilterOperator = "gt"
	// GreaterOrEqual matches values greater than or equal to the specified value.
	GreaterOrEqual FilterOperator = "gte"
	// Less matches values less than the specified value.
	Less FilterOperator = "lt"
	// LessOrEqual matches values less than or equal to the specified value.
	LessOrEqual FilterOperator = "lte"
	// Like matches values matching a pattern.
	Like FilterOperator = "like"
	// In matches values in a list.
	In FilterOperator = "in"
	// NotIn matches values not in a list.
	NotIn FilterOperator = "nin"
	// Contains matches values containing a substring.
	Contains FilterOperator = "contains"
	// NotContains matches values not containing a substring.
	NotContains FilterOperator = "ncontains"
	// Is checks for identity (e.g., IS NULL).
	Is FilterOperator = "is"
	// IsNull checks if value is NULL.
	IsNull FilterOperator = "isnull"
	// IsNot checks for non-identity.
	IsNot FilterOperator = "isn"
	// IsNotNull checks if value is NOT NULL.
	IsNotNull FilterOperator = "isnnull"

	// Default is used when no operator is specified and the value is already defined with a sub-query.
	Default FilterOperator = "default"
)

// FilterChainingKey represents the logical operator to chain filters.
type FilterChainingKey string

const (
	// And represents the AND logical operator.
	And FilterChainingKey = "AND"
	// Or represents the OR logical operator.
	Or FilterChainingKey = "OR"
)

// Filter represents a single filter criteria.
type Filter struct {
	Module                            string
	IsGroupOpen                       bool
	GroupOpenQty                      int
	Field                             FilterField
	Operator                          FilterOperator
	Value                             FilterValue
	IsGroupClose                      bool
	GroupCloseQty                     int
	ChainingKey                       FilterChainingKey
	OverridePreviousFilterChainingKey FilterChainingKey
}

// Filters represents a collection of filters.
type Filters []Filter

// IsZero returns true if the filters collection is empty.
func (f Filters) IsZero() bool {
	return len(f) == 0
}

// FilterBy creates a new Filters collection with a single filter.
func FilterBy(name string, operator FilterOperator, value any) Filters {
	return Filters{
		{
			Field:    FilterField(name),
			Operator: operator,
			Value:    value,
		},
	}
}

// Or adds a new filter combined with OR logic.
func (f Filters) Or(field string, operator FilterOperator, value any) Filters {
	if f.IsZero() {
		return Filters{{Field: FilterField(field), Operator: operator, Value: value}}
	}

	f[len(f)-1].ChainingKey = Or

	return append(f, Filter{
		Field:    FilterField(field),
		Operator: operator,
		Value:    value,
	})
}

// And adds a new filter combined with AND logic.
func (f Filters) And(field string, operator FilterOperator, value any) Filters {
	if f.IsZero() {
		return Filters{{Field: FilterField(field), Operator: operator, Value: value}}
	}

	f[len(f)-1].ChainingKey = And

	return append(f, Filter{
		Field:    FilterField(field),
		Operator: operator,
		Value:    value,
	})
}

// AndGroup adds a group of filters combined with AND logic.
func (f Filters) AndGroup(filters ...Filter) Filters {
	if len(filters) == 0 {
		return f
	}

	if len(f) > 0 {
		f[len(f)-1].ChainingKey = And
	}

	filters[0].IsGroupOpen = true
	filters[len(filters)-1].IsGroupClose = true

	f = append(f, filters...)

	return f
}

// OrGroup adds a group of filters combined with OR logic.
func (f Filters) OrGroup(filters ...Filter) Filters {
	if len(filters) == 0 {
		return f
	}

	if len(f) > 0 {
		f[len(f)-1].ChainingKey = Or
	}

	filters[0].IsGroupOpen = true
	filters[len(filters)-1].IsGroupClose = true

	f = append(f, filters...)

	return f
}
