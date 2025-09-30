package performance

import (
	"testing"

	"github.com/go-interpreter/internal/scanner"
)

func BenchmarkScannerBinarySimple(b *testing.B) {
	programText := "(1+2)"
	for b.Loop() {
		tokenScanner := scanner.NewTokenScanner(programText)
		_ = tokenScanner.ScanTokens()
	}
}

func BenchmarkScannerBinaryComplex(b *testing.B) {
	programText := "(1+2+(3+4))"
	for b.Loop() {
		tokenScanner := scanner.NewTokenScanner(programText)
		_ = tokenScanner.ScanTokens()
	}
}

func BenchmarkScannerStatements(b *testing.B) {
	programText := "print 10 + 20;"
	for b.Loop() {
		tokenScanner := scanner.NewTokenScanner(programText)
		_ = tokenScanner.ScanTokens()
	}

}
