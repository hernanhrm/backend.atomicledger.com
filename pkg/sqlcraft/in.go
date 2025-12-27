package sqlcraft

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

// In builds an IN clause.
func In(value any, initialArgCount int) Result {
	if value == nil {
		return Result{}
	}

	builder := bytes.Buffer{}
	builder.WriteString("(")

	// uses reflection to handle different types.
	valSlice := reflect.ValueOf(value)
	if valSlice.Kind() == reflect.Slice {
		if valSlice.Len() == 0 {
			return Result{}
		}

		args := make([]any, 0, valSlice.Len())
		for i := range valSlice.Len() {
			builder.WriteString("$")
			builder.WriteString(strconv.Itoa(initialArgCount + i))
			builder.WriteString(", ")

			args = append(args, valSlice.Index(i).Interface())
		}

		if valSlice.Len() > 0 {
			builder.Truncate(builder.Len() - 2)
		}

		builder.WriteString(")")

		return Result{
			SQL:  builder.String(),
			Args: args,
		}
	}

	str, ok := value.(string)
	if !ok {
		return Result{}
	}

	stringValues := strings.Split(str, ",")
	args := make([]any, 0, len(stringValues))
	for i, v := range stringValues {
		builder.WriteString("$")
		builder.WriteString(strconv.Itoa(initialArgCount + i))
		builder.WriteString(", ")

		args = append(args, v)
	}

	builder.Truncate(builder.Len() - 2)
	builder.WriteString(")")

	return Result{
		SQL:  builder.String(),
		Args: args,
	}
}
