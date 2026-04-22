package rules

import (
	"ConfigScanner/pkg/models"
	"strings"
)

// BindAllRule структура для правила "прослушка всех интерфейсов"
// Это правило ищет опасную практику, когда сервер слушает на всех сетевых интерфейсах (0.0.0.0)
// Сервер становится доступным извне (из интернета)
// Если нет firewall, любой может подключиться
type BindAllRule struct{}

func (r BindAllRule) Name() string {
	return "bind-all"
}

func (r BindAllRule) Check(config map[string]interface{}, path string) []models.Problem {
	var problems []models.Problem

	for key, value := range config {
		currentPath := joinPath(path, key)

		// Проверяем ключи, связанные с хостом
		if isHostKey(key) {
			if str, ok := value.(string); ok && str == "0.0.0.0" {
				problems = append(problems, models.Problem{
					Rule:           r.Name(),
					LevelProblem:   models.Medium, // В разработке это допустимо локальный компьютер, продакшене это опасно.
					Path:           currentPath,
					ParseV:         "Сервер слушает на всех интерфейсах (0.0.0.0)",
					Recommendation: "Ограничьте доступ, используйте localhost или конкретный IP-адрес",
				})
			}
		}

		// Рекурсивно проверяем вложенные структуры
		if nested, ok := value.(map[string]interface{}); ok {
			problems = append(problems, r.Check(nested, currentPath)...)
		}
	}

	return problems
}

// Проверяем, является ли текущий ключ ключом, связанным с хостом, адресом
// isHostKey проверяет ключи: "host", "bind", "address", "listen"
func isHostKey(key string) bool {

	lower := strings.ToLower(key)
	return lower == "host" || lower == "bind" || lower == "address" || lower == "listen"
}
