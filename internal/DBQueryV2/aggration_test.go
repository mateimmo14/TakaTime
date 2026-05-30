package dbqueryv2

import (
	"testing"
	"time"
)

func TestCalculateTrueDuration(t *testing.T) {
	base := time.Unix(1000, 0)

	tests := []struct {
		name     string
		logs     []LogTime
		expected float64
	}{
		{
			name: "single log",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
			},
			expected: 60,
		},
		{
			name: "overlapping logs",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
				{Timestamp: base.Add(30 * time.Second), Duration: 60},
			},
			expected: 90,
		},
		{
			name: "non-overlapping logs",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
				{Timestamp: base.Add(120 * time.Second), Duration: 60},
			},
			expected: 120,
		},
		{
			name:     "empty logs",
			logs:     []LogTime{},
			expected: 0,
		},
		{
			name: "contained overlap",
			logs: []LogTime{
				{Timestamp: base, Duration: 120},
				{Timestamp: base.Add(-30 * time.Second), Duration: 30},
			},
			expected: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTrueDuration(tt.logs)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "minutes only",
			seconds:  300,
			expected: "5m",
		},
		{
			name:     "hours and minutes",
			seconds:  3660,
			expected: "1h 1m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.seconds)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
