package helper_test

import (
	"testing"

	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
	"github.com/TheAmirhosssein/cool-password-manage/internal/utils/helper"
	"github.com/stretchr/testify/assert"
)

func TestMakeSearchQuery(t *testing.T) {
	tests := []struct {
		name   string
		q      types.NullString
		fields []string
		want   string
	}{
		{
			name: "valid single field",
			q: types.NullString{
				String: "john",
				Valid:  true,
			},
			fields: []string{"name"},
			want:   "AND (name ILIKE '%john%')",
		},
		{
			name: "valid multiple fields",
			q: types.NullString{
				String: "john",
				Valid:  true,
			},
			fields: []string{"name", "email"},
			want:   "AND (name ILIKE '%john%' OR email ILIKE '%john%')",
		},
		{
			name: "invalid null string",
			q: types.NullString{
				Valid: false,
			},
			fields: []string{"name"},
			want:   "",
		},
		{
			name: "empty string",
			q: types.NullString{
				String: "",
				Valid:  true,
			},
			fields: []string{"name"},
			want:   "",
		},
		{
			name: "no fields",
			q: types.NullString{
				String: "john",
				Valid:  true,
			},
			fields: nil,
			want:   "",
		},
		{
			name: "escape single quote",
			q: types.NullString{
				String: "o'connor",
				Valid:  true,
			},
			fields: []string{"name"},
			want:   "AND (name ILIKE '%o''connor%')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.MakeSearchQuery(tt.q, tt.fields)
			assert.Equal(t, tt.want, got)
		})
	}
}
