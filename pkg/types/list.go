// Package types provides generic data types and utilities.
package types

type (
	// ListIterator is a function type for iterating over a List.
	// It receives the index and the value, and returns a boolean indicating whether to proceed/select.
	ListIterator[T any] = func(k int, v T) bool
)

// List is a generic slice type with utility methods.
type List[T any] []T

// IsEmpty returns true if the list is empty.
func (l List[T]) IsEmpty() bool {
	return len(l) == 0
}

// Filter returns a new List containing only the elements for which the iterator returns true.
func (l List[T]) Filter(iter ListIterator[T]) []T {
	res := make([]T, 0, len(l))
	for k, v := range l {
		if iter(k, v) {
			res = append(res, v)
		}
	}

	return res
}

// Find returns the first element for which the iterator returns true, and a boolean indicating success.
func (l List[T]) Find(iter ListIterator[T]) (T, bool) {
	for k, v := range l {
		if iter(k, v) {
			return v, true
		}
	}

	return *new(T), false
}

// MustFind returns the first element for which the iterator returns true, or the zero value of T if found.
// Note: This name implies it might panic, but currently it returns zero value.
func (l List[T]) MustFind(iter ListIterator[T]) T {
	for k, v := range l {
		if iter(k, v) {
			return v
		}
	}

	return *new(T)
}
