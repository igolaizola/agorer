package isbn

import (
	"context"
	"testing"
	"time"
)

func TestHyphenate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "v for vendetta",
			input: "9781779511195",
			want:  "978-1-77951-119-5",
		},
		{
			name:  "binti",
			input: "9788494795886",
			want:  "978-84-947958-8-6",
		},
		{
			name:  "el último minuto",
			input: "9788418054525",
			want:  "978-84-18054-52-5",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := New("", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if !Valid(tt.input) {
				t.Errorf("invalid isbn: %s", tt.input)
			}
			got, err := client.Hyphenate(ctx, tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Hyphenate() = %v, want %v", got, tt.want)
			}
		})
	}
}
