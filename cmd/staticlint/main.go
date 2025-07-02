/*
Package main запускает multichecker, состоящий из набора статических анализаторов кода Go.

Включает:

  - Стандартные анализаторы из пакета golang.org/x/tools/go/analysis/passes: printf, shadow, structtag

  - Все анализаторы класса SA из пакета staticcheck.io.

  - Один дополнительный анализатор из staticcheck.io (S1005).

  - github.com/Antonboom/errname — проверяет корректность наименования переменных ошибок.

  - github.com/tdakkota/asciicheck — проверяет, что идентификаторы написаны в ASCII.

  - Пользовательский анализатор noexit, запрещающий использование прямого вызова os.Exit
    в функции main пакета main.

Запуск:

	go run ./cmd/staticlint ./...
*/
package main

import (
	errname "github.com/Antonboom/errname/pkg/analyzer"
	"github.com/invinciblewest/metrics/internal/analyzer/noexit"
	"github.com/tdakkota/asciicheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		noexit.Analyzer,
		errname.New(),
		asciicheck.NewAnalyzer(),
	)

	for _, a := range staticcheck.Analyzers {
		if len(a.Analyzer.Name) >= 2 && a.Analyzer.Name[:2] == "SA" || a.Analyzer.Name == "S1005" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	multichecker.Main(analyzers...)
}
