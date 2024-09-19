package chatgpt

import (
	"testing"
)

func Test_removeHeaders(t *testing.T) {
	response := "<|start_header_id|>assistant<|end_header_id|>\\\\{\"title\": \"Title\"}"
	expected := "{\"title\": \"Title\"}"
	got := removeHeaders(response)

	if expected != got {
		t.Errorf("unexpected response, expected: %q, got: %q", expected, got)
	}
}
