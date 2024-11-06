package models

import (
	"testing"

	"github.com/tneuqole/snippetbox/internal/models/assert"
)

func TestUserModelExists(t *testing.T) {
	// skip if -short flag is provided
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name string
		ID   int
		want bool
	}{
		{
			name: "Valid ID",
			ID:   1,
			want: true,
		},
		{
			name: "Zero ID",
			ID:   0,
			want: false,
		},
		{
			name: "Non-existent ID",
			ID:   2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(tt.ID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
