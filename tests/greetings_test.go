package greetings

import (
	"testing"
)

func TestHello(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal name", "Chris", "Hello, Chris"},
		{"empty string", "", "Hello, World"},
		{"name with space", "Chris Evans", "Hello, Chris Evans"},
		{"special characters", "O'Connor", "Hello, O'Connor"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hello(tt.input)
			if got != tt.want {
				t.Errorf("got %q want %q", got, tt.want)
			}
		})
	}
}
