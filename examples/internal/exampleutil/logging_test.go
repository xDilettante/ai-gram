package exampleutil

import "testing"

func TestMaskInt64(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{name: "zero", value: 0, want: "0"},
		{name: "short positive", value: 123456, want: "***"},
		{name: "short negative", value: -12345, want: "***"},
		{name: "long positive", value: 123456789, want: "123***789"},
		{name: "long negative", value: -1001234567890, want: "-10***890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskInt64(tt.value); got != tt.want {
				t.Fatalf("MaskInt64(%d) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}
