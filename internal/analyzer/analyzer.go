package analyzer

import (
	"ConfigScanner/internal/rules"
	"ConfigScanner/pkg/models"
)

// Analyzer анализирует конфигурацию по набору правил
// Структура хранит слайс правил, каждое из которых реализует интерфейс Rule
type Analyzer struct {
	rules []rules.Rule
}

// NewAnalyzer создает анализатор с правилами из ТЗ
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		rules: []rules.Rule{
			&rules.DebugLogRule{},
			&rules.PlainPasswordRule{},
			&rules.BindAllRule{},
			&rules.TLSDisabledRule{},
			&rules.WeakAlgorithmRule{},
		},
	}
}

// Analyze проверяет конфигурацию и возвращает список проблем
// config: распарсенный конфиг в виде map[string]interface{}
// Возвращает: []models.Problem: слайс найденных проблем (каждая содержит уровень,
// сообщение и рекомендацию). Если проблем нет, возвращается пустой слайс.
// Алгоритм работы:
//  1. Создается пустой слайс allProblems
//  2. Для каждого правила из a.rules вызывается метод Check()
//  3. Check() возвращает слайс проблем для этого правила (может быть пустым)
//  4. Оператор ... распаковывает элементы problems и добавляет их в allProblems
//  5. В конце возвращается общий слайс всех найденных проблем
func (a *Analyzer) Analyze(config map[string]interface{}) []models.Problem {
	var allProblems []models.Problem

	for _, rule := range a.rules {
		problems := rule.Check(config, "")
		allProblems = append(allProblems, problems...) //Три точки (...) в Go — это оператор распаковки для slices.
		// Можно написать так, но код длинее:
		/*for _, problem := range problems {
			allProblems = append(allProblems, problem)
		}*/
	}

	return allProblems
}
