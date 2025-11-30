package utils

import (
	"testing"
	"time"
)

func TestRelativeAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		date     string
		expected string
	}{
		{
			name:     "just now",
			date:     now.Add(-30 * time.Second).Format("2006-01-02 15:04:05"),
			expected: "30 seconds ago",
		},
		{
			name:     "minutes ago",
			date:     now.Add(-10 * time.Minute).Format("2006-01-02 15:04:05"),
			expected: "10 minutes ago",
		},
		{
			name:     "hours ago",
			date:     now.Add(-5 * time.Hour).Format("2006-01-02 15:04:05"),
			expected: "5 hours ago",
		},
		{
			name:     "days ago",
			date:     now.Add(-48 * time.Hour).Format("2006-01-02 15:04:05"),
			expected: "2 days ago",
		},
		{
			name:     "invalid date",
			date:     "invalid-date",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generation{Date: tt.date}
			got := g.RelativeAge()
			if got != tt.expected {
				t.Errorf("RelativeAge() = %q, want %q", got, tt.expected)
			}
		})
	}
}
