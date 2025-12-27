// Package sqlcraft provides a fluent SQL query builder for PostgreSQL.
package sqlcraft

import "errors"

var (
	// ErrEmptyValues is returned when no values are provided for the INSERT query.
	ErrEmptyValues = errors.New("empty values in query")
	// ErrEmptyColumns is returned when no columns are provided for the INSERT/SELECT query.
	ErrEmptyColumns = errors.New("empty columns in query")
	// ErrMissMatchValues is returned when the number of values does not match the number of columns.
	ErrMissMatchValues = errors.New("miss match values for given columns")
	// ErrInvalidOperator is returned when an invalid operator is encountered.
	ErrInvalidOperator = errors.New("invalid dafi operator")
	// ErrInvalidFieldName is returned when an invalid field name is encountered.
	ErrInvalidFieldName = errors.New("invalid field name")
)
