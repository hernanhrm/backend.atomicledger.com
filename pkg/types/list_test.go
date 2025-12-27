// Package types_test contains tests for the types package.
package types_test

import (
	"testing"

	"backend.atomicledger.com/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestList_IsEmpty(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		var l types.List[int]
		assert.True(t, l.IsEmpty())
	})

	t.Run("non-empty list", func(t *testing.T) {
		l := types.List[int]{1, 2, 3}
		assert.False(t, l.IsEmpty())
	})
}

func TestList_Filter(t *testing.T) {
	l := types.List[int]{1, 2, 3, 4, 5, 6}
	filtered := l.Filter(func(_ int, v int) bool {
		return v%2 == 0
	})
	assert.Equal(t, []int{2, 4, 6}, filtered)

	t.Run("filter everything", func(t *testing.T) {
		filtered := l.Filter(func(_ int, _ int) bool {
			return false
		})
		assert.Empty(t, filtered)
	})

	t.Run("filter nothing", func(t *testing.T) {
		filtered := l.Filter(func(_ int, _ int) bool {
			return true
		})
		assert.Equal(t, []int(l), filtered)
	})
}

func TestList_Find(t *testing.T) {
	l := types.List[string]{"apple", "banana", "cherry"}

	t.Run("found", func(t *testing.T) {
		val, found := l.Find(func(_ int, v string) bool {
			return v == "banana"
		})
		assert.True(t, found)
		assert.Equal(t, "banana", val)
	})

	t.Run("not found", func(t *testing.T) {
		val, found := l.Find(func(_ int, v string) bool {
			return v == "dragonfruit"
		})
		assert.False(t, found)
		assert.Equal(t, "", val)
	})
}

func TestList_MustFind(t *testing.T) {
	l := types.List[int]{10, 20, 30}

	t.Run("found", func(t *testing.T) {
		val := l.MustFind(func(_ int, v int) bool {
			return v == 20
		})
		assert.Equal(t, 20, val)
	})

	t.Run("not found", func(t *testing.T) {
		val := l.MustFind(func(_ int, v int) bool {
			return v == 99
		})
		assert.Equal(t, 0, val)
	})
}
