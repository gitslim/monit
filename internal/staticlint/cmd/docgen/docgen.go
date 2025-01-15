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
		panic(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, _ = file.WriteString(`// Команда staticlint запускает multichecker для статического анализа кода.
//
// Для запуска multichecker используйте команду:
//
// go run cmd/staticlint/main.go
//
// Используются следующие анализаторы:
//
`)

	for _, analyzer := range analyzers {
		_, _ = file.WriteString("//\n// " + staticlint.AnalyzerDescription(analyzer) + "\n")
	}

	_, _ = file.WriteString("package main\n")

	fmt.Println("Staticlint documentation generated successfully!")
}
