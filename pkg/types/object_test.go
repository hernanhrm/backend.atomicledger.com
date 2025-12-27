// Package types_test contains tests for the types package.
package types_test

import (
	"encoding/json"
	"testing"

	"backend.atomicledger.com/pkg/types"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestObject_NewObject(t *testing.T) {
	data := TestData{Name: "John", Age: 30}
	obj := types.NewObject(data)
	assert.Equal(t, data, obj.Data())
}

func TestObject_Value(t *testing.T) {
	data := TestData{Name: "John", Age: 30}
	obj := types.NewObject(data)

	val, err := obj.Value()
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"John","age":30}`, string(val.([]byte)))
}

func TestObject_Scan(t *testing.T) {
	t.Run("scan bytes", func(t *testing.T) {
		jsonBytes := []byte(`{"name":"Jane","age":25}`)
		var obj types.Object[TestData]
		err := obj.Scan(jsonBytes)
		assert.NoError(t, err)
		assert.Equal(t, "Jane", obj.Data().Name)
		assert.Equal(t, 25, obj.Data().Age)
	})

	t.Run("scan string", func(t *testing.T) {
		jsonStr := `{"name":"Bob","age":40}`
		var obj types.Object[TestData]
		err := obj.Scan(jsonStr)
		assert.NoError(t, err)
		assert.Equal(t, "Bob", obj.Data().Name)
		assert.Equal(t, 40, obj.Data().Age)
	})

	t.Run("scan invalid type", func(t *testing.T) {
		var obj types.Object[TestData]
		err := obj.Scan(123)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal JSONB value")
	})

	t.Run("scan invalid json", func(t *testing.T) {
		var obj types.Object[TestData]
		err := obj.Scan([]byte(`{invalid json}`))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal Object data")
	})
}

func TestObject_MarshalJSON(t *testing.T) {
	data := TestData{Name: "Alice", Age: 28}
	obj := types.NewObject(data)

	bytes, err := json.Marshal(obj)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"name":"Alice","age":28}`, string(bytes))
}

func TestObject_UnmarshalJSON(t *testing.T) {
	jsonBytes := []byte(`{"name":"Charlie","age":35}`)
	var obj types.Object[TestData]
	err := json.Unmarshal(jsonBytes, &obj)

	assert.NoError(t, err)
	assert.Equal(t, "Charlie", obj.Data().Name)
	assert.Equal(t, 35, obj.Data().Age)
}
