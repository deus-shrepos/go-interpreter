package scanner

import (
	"testing"

	"github.com/go-interpreter/internal/scanner"
)

func TestNewTokenScanner(t *testing.T) {
	t.Run("TestNewTokenScanner", func(t *testing.T) {
		programText := "(1+2)"
		tokenScanner := scanner.NewTokenScanner(programText)
		if tokenScanner.Source != programText {
			t.Errorf("Expected Source to be %s, got %s", programText, tokenScanner.Source)
		}
		if tokenScanner.Current != 0 {
			t.Errorf("Expected Current to be 0, got %d", tokenScanner.Current)
		}
		if tokenScanner.Line != 0 {
			t.Errorf("Expected Line to be 0, got %d", tokenScanner.Line)
		}
		if tokenScanner.Start != 0 {
			t.Errorf("Expected Start to be 0, got %d", tokenScanner.Start)
		}
	})
}
