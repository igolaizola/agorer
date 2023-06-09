package sinli

import (
	"reflect"
	"testing"
)

// Example usage
type Bar struct {
	_       struct{} `sinli:"order=1,length=1,fixed=B"`
	Text    string   `sinli:"order=2,length=10"`
	Number  int      `sinli:"order=3,length=5"`
	Boolean bool     `sinli:"order=4,length=1"`
}

type Foo struct {
	Header Bar   `sinli:"order=1"`
	Body   []Bar `sinli:"order=2"`
}

func TestSinliMarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{
			name: "struct",
			input: Bar{
				Text:    "Hello",
				Number:  42,
				Boolean: true,
			},
			want: []byte("BHello     00042S\r\n"),
		},
		{
			name: "pointer",
			input: &Bar{
				Text:    "Hello",
				Number:  -42,
				Boolean: true,
			},
			want: []byte("BHello     -0042S\r\n"),
		},
		{
			name: "slice",
			input: []Bar{
				{
					Text:    "Hello",
					Number:  42,
					Boolean: true,
				},
				{
					Text:    "World",
					Number:  -43,
					Boolean: false,
				},
			},
			want: []byte("BHello     00042S\r\nBWorld     -0043N\r\n"),
		},
		{
			name: "array",
			input: [2]Bar{
				{
					Text:    "Hello",
					Number:  42,
					Boolean: true,
				},
				{
					Text:    "World",
					Number:  -43,
					Boolean: false,
				},
			},
			want: []byte("BHello     00042S\r\nBWorld     -0043N\r\n"),
		},
		{
			name: "nested",
			input: Foo{
				Header: Bar{
					Text:    "Header",
					Number:  1,
					Boolean: true,
				},
				Body: []Bar{
					{
						Text:    "Hello",
						Number:  42,
						Boolean: true,
					},
					{
						Text:    "World",
						Number:  -43,
						Boolean: false,
					},
				},
			},
			want: []byte("BHeader    00001S\r\nBHello     00042S\r\nBWorld     -0043N\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			got, err := Marshal(tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil, got error: %s", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("want:\n'%s'\ngot:\n'%s'\n", tt.want, got)
			}
		})
	}
}
