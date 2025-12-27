package sqlcraft

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert_ToSql(t *testing.T) {
	type fields struct {
		table            string
		columns          []string
		returningColumns []string
		values           []any
	}
	tests := []struct {
		name    string
		query   InsertQuery
		fields  fields
		want    Result
		wantErr bool
	}{
		{
			name:  "standard insert",
			query: InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash"),
			want: Result{
				SQL:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			wantErr: false,
		},
		{
			name:  "standard insert with returning",
			query: InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash").Returning("id", "created_at"),
			want: Result{
				SQL:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash"},
			},
			wantErr: false,
		},
		{
			name: "standard insert with returning and multiple row values",
			fields: fields{
				table:            "users",
				columns:          []string{"first_name", "last_name", "email", "password"},
				returningColumns: []string{"id", "created_at"},
				values:           []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash", "Brownie", nil, "brownie@gmail.com", "secrethash"},
			},
			query: InsertInto("users").
				WithColumns("first_name", "last_name", "email", "password").
				WithValues("Hernan", nil, "hernan_rm@outlook.es", "secrethash").
				WithValues("Brownie", nil, "brownie@gmail.com", "secrethash").
				Returning("id", "created_at"),
			want: Result{
				SQL:  "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4), ($5, $6, $7, $8) RETURNING id, created_at",
				Args: []any{"Hernan", nil, "hernan_rm@outlook.es", "secrethash", "Brownie", nil, "brownie@gmail.com", "secrethash"},
			},
			wantErr: false,
		},
		{
			name:    "error empty values",
			query:   InsertInto("users").WithColumns("first_name", "last_name", "email", "password").Returning("id", "created_at"),
			want:    Result{},
			wantErr: true,
		},
		{
			name:    "missmatch values",
			query:   InsertInto("users").WithColumns("first_name", "last_name", "email", "password").WithValues("hernan").Returning("id", "created_at"),
			want:    Result{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.query.ToSQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
