package rules

import (
	"ConfigScanner/pkg/models"
	"strings"
)

// TLSDisabledRule структура для правила "отключенный TLS"
// Это правило ищет места, где TLS шифрование отключено
// Чем это опасно:
// Все данные передаются в открытом виде
// Злоумышленник может перехватить пароли, токены, личные данные
// Пример опасного конфига:
//
//	server:
//	  tls:
//	    enabled: false  # ← Вот это ищем!
type TLSDisabledRule struct{}

func (r TLSDisabledRule) Name() string {
	return "tls-disabled"
}

func (r TLSDisabledRule) Check(config map[string]interface{}, path string) []models.Problem {
	var problems []models.Problem

	for key, value := range config {
		currentPath := joinPath(path, key)

		// Проверяем ключ enabled/enable
		if isEnabledKey(key) {
			if enabled, ok := value.(bool); ok && !enabled {
				// Проверяем, относится ли это к TLS
				if isTLSContext(path, key) {
					problems = append(problems, models.Problem{
						Rule:           r.Name(),
						LevelProblem:   models.High, //Отключенный TLS в production - это критическая уязвимость!
						Path:           currentPath,
						ParseV:         "TLS отключен",
						Recommendation: "Включите TLS для безопасного соединения",
					})
				}
			}
		}

		// Рекурсивно проверяем вложенные структуры
		if nested, ok := value.(map[string]interface{}); ok {
			problems = append(problems, r.Check(nested, currentPath)...)
		}
	}

	return problems
}

// isEnabledKey проверяет, является ли ключ ключом включения/выключения
func isEnabledKey(key string) bool {
	lower := strings.ToLower(key)
	return lower == "enabled" || lower == "enable"
}

// isTLSContext проверяет, относится ли поле enabled к TLS
func isTLSContext(path, key string) bool {
	fullPath := strings.ToLower(path + "." + key)
	return strings.Contains(fullPath, "tls") ||
		strings.Contains(fullPath, "https")
}
