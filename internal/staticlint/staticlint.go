// Package staticlint предоставляет анализаторы для статического анализа кода.
package staticlint

import (
	"fmt"
	"strings"

	"github.com/gitslim/monit/internal/staticlint/analyzers/exitcheck"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func AnalyzerDescription(a *analysis.Analyzer) string {
	return fmt.Sprintf("%s: %s", a.Name, strings.Split(a.Doc, "\n")[0])
}

// AllAnalyzers возвращает список всех анализаторов.
func AllAnalyzers() []*analysis.Analyzer {
	// добавлние стандартных анализаторов.
	all := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}

	// добавление анализатора errcheck.
	all = append(all, errcheck.Analyzer)

	// добавление кастомного анализатора exitcheck.
	all = append(all, exitcheck.Analyzer)

	// добавление всех S-анализаторов simple.
	for _, v := range simple.Analyzers {
		all = append(all, v.Analyzer)
	}

	// добавление всех SA-анализаторов staticcheck.
	for _, v := range staticcheck.Analyzers {
		all = append(all, v.Analyzer)
	}

	// добавление всех ST-анализаторов stylecheck.
	for _, v := range stylecheck.Analyzers {
		all = append(all, v.Analyzer)
	}

	return all
}
