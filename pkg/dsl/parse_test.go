package dsl

import (
	"path/filepath"
	"testing"
)

func TestParseFile_Examples(t *testing.T) {
	examples := []struct {
		file      string
		wantValid bool
	}{
		{"example-01-minimal.yaml", true},
		{"example-02-full-features.yaml", true},
		{"example-03-embedded.yaml", true},
		{"example-04-polymorphic-go.yaml", true},
	}

	for _, tt := range examples {
		t.Run(tt.file, func(t *testing.T) {
			path := filepath.Join("..", "..", "docs", "gp", "08-dsl-examples", tt.file)
			spec, err := ParseFile(path)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			errs := Validate(spec)
			if tt.wantValid && len(errs) > 0 {
				t.Errorf("expected valid, got errors: %v", errs)
			}
		})
	}
}
