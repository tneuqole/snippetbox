package main

import (
	"testing"
	"time"

	"github.com/tneuqole/snippetbox/internal/models/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		x    time.Time
		want string
	}{
		{
			name: "UTC",
			x:    time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			x:    time.Time{},
			want: "",
		},
		{
			name: "EST",
			x:    time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("EST", -5*60*60)),
			want: "17 Mar 2024 at 15:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := humanDate(tt.x)
			assert.Equal(t, d, tt.want)
		})
	}
}
