// Package ternary provides utility functions for ternary operations.
package ternary

// If returns trueValue if condition is true, otherwise returns falseValue.
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
