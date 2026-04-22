package rule

import "ConfigScanner/pkg/models"

// Rule интерфейс для всех правил проверки
type Rule interface {
	// Name возвращает название правила
	Name() string

	// Check проверяет конфиг и возвращает проблемы
	Check(config map[string]interface{}, path string) []models.Problem
}
