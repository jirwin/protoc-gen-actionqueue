package pgactionqueue

import "testing"

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"my-entity-queue", "MyEntityQueue"},
		{"book-queue", "BookQueue"},
		{"snake_case", "SnakeCase"},
		{"single", "Single"},
		{"", ""},
		{"my-mixed_case", "MyMixedCase"},
		{"ALL-CAPS", "AllCaps"},
		{"a-b-c", "ABC"},
		{"queue-v2-beta", "QueueV2Beta"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
