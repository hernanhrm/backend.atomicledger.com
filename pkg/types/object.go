// Package types provides generic data types and utilities.
package types

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/samber/oops"
)

// Object is a generic type that implements sql.Scanner and driver.Valuer interfaces.
// for handling single objects of any type as database values.
// NOTE: always ensure that you use Object as a pointer, if not, the data will appear as a field in the json.
type Object[T any] struct {
	data T
}

// NewObject creates a new Object instance with the provided data.
func NewObject[T any](data T) Object[T] {
	return Object[T]{data: data}
}

// Data returns the underlying data of the Object.
func (op *Object[T]) Data() T {
	return op.data
}

// Value return json value, implement driver.Valuer interface.
func (op Object[T]) Value() (driver.Value, error) {
	bytes, err := json.Marshal(op.data)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to marshal Object data")
	}
	return bytes, nil
}

// Scan scan value into JSONType[T], implements sql.Scanner interface.
func (op *Object[T]) Scan(value any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return oops.Errorf("failed to unmarshal JSONB value: %v", value)
	}

	if err := json.Unmarshal(bytes, &op.data); err != nil {
		return oops.Wrapf(err, "failed to unmarshal Object data")
	}
	return nil
}

// MarshalJSON to output non base64 encoded []byte.
func (op Object[T]) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(op.data)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to marshal Object data")
	}
	return bytes, nil
}

// UnmarshalJSON to deserialize []byte.
func (op *Object[T]) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &op.data); err != nil {
		return oops.Wrapf(err, "failed to unmarshal Object data")
	}
	return nil
}
