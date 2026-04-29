package wdkgo

import (
	"path/filepath"
	"testing"

	"solod.dev/internal/compiler"
)

func TestTranslate(t *testing.T) {
	srcDir := filepath.Join(".", "hyperdbg-driver")
	outDir := filepath.Join(".", "hyperdbg-driver", "gen")

	if err := compiler.Translate(srcDir, outDir); err != nil {
		t.Fatalf("translate failed: %v", err)
	}
}
