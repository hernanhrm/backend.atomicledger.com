package ternary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	t.Run("should return trueValue when condition is true", func(t *testing.T) {
		result := If(true, "yes", "no")
		assert.Equal(t, "yes", result)
	})

	t.Run("should return falseValue when condition is false", func(t *testing.T) {
		result := If(false, "yes", "no")
		assert.Equal(t, "no", result)
	})

	t.Run("should work with integers", func(t *testing.T) {
		result := If(true, 1, 0)
		assert.Equal(t, 1, result)
	})

	t.Run("should work with structs", func(t *testing.T) {
		type user struct{ name string }
		u1 := user{name: "Alice"}
		u2 := user{name: "Bob"}
		result := If(false, u1, u2)
		assert.Equal(t, u2, result)
	})
}
