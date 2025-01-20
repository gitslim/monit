package main

//go:generate go run ../../internal/staticlint/cmd/docgen/docgen.go

import (
	"github.com/gitslim/monit/internal/staticlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	analyzers := staticlint.AllAnalyzers()
	// fmt.Println("Loading analyzers...")
	// for _, analyzer := range analyzers {
	// 	fmt.Println(staticlint.AnalyzerDescription(analyzer))
	// }
	// fmt.Println("Running multichecker...")
	// fmt.Println()
	multichecker.Main(analyzers...)
}
