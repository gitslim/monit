package main

import (
	"fmt"
	"os"

	"github.com/gitslim/monit/internal/staticlint"
)

func main() {
	analyzers := staticlint.AllAnalyzers()

	file, err := os.Create("doc.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create documentation file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, _ = file.WriteString(`// Command staticlint запускает multichecker для статического анализа кода.
//
// Для запуска multichecker используйте команду:
//
// go run cmd/staticlint/main.go
//
// Используются следующие анализаторы:
//
`)

	for _, analyzer := range analyzers {
		_, _ = file.WriteString("// " + staticlint.AnalyzerDescription(analyzer) + "\n")
	}

	_, _ = file.WriteString("package main\n")

	fmt.Println("Staticlint documentation generated successfully!")
}
