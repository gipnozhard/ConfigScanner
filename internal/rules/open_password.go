package rules

import (
	"ConfigScanner/pkg/models"
	"strings"
)

// PlainPasswordRule структура для правила "пароль в открытом виде"
// Это правило ищет пароли, секреты и токены, которые хранятся в конфиге в незашифрованном виде
type PlainPasswordRule struct{}

func (r PlainPasswordRule) Name() string {
	return "plain-password"
}

func (r PlainPasswordRule) Check(config map[string]interface{}, path string) []models.Problem {
	var problems []models.Problem

	for key, value := range config {
		currentPath := joinPath(path, key)

		// Проверяем ключи, похожие на пароли
		if isPasswordKey(key) {
			if isSafeKey(key) {
				continue // не считаем это проблемой
			}

			if str, ok := value.(string); ok && str != "" && !isEnvReference(str) {
				problems = append(problems, models.Problem{
					Rule:           r.Name(),
					LevelProblem:   models.High,
					Path:           currentPath,
					ParseV:         "Обнаружен пароль/секрет в открытом виде",
					Recommendation: "Используйте переменные окружения",
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

// isPasswordKey проверяет, является ли ключ ключом уровня пароля
func isPasswordKey(key string) bool {
	lower := strings.ToLower(key)
	return strings.Contains(lower, "password") ||
		strings.Contains(lower, "passwd") ||
		strings.Contains(lower, "secret") ||
		strings.Contains(lower, "token") ||
		// Для ключей (key) делаем дополнительную проверку:
		// Ищем "key", НО исключаем "public" и "private", потому что:
		// public_key, publicKey - это открытые ключи (не секрет)
		// private_key, privateKey - хотя это секрет, но обычно хранится в отдельном файле
		//     и в конфиге указывается только путь к файлу
		strings.Contains(lower, "key") && !strings.Contains(lower, "public") && !strings.Contains(lower, "private")
}

// isEnvReference проверяет, является ли значение ссылкой на переменную окружения
// value: строка со значением из конфига (например "DB_PASSWORD", "${DB_PASS}", "$PASSWORD")
// Возвращает:
// true если значение выглядит как ссылка на переменную окружения
// false если это реальное значение (например "admin123")
// Это нужно, чтобы не срабатывать на безопасные варианты хранения секретов:
// password: "${DB_PASSWORD}" - переменная окружения (безопасно)
// secret: "$API_SECRET" - переменная окружения (безопасно)
// token: "ENV:JWT_TOKEN" - специальный маркер (безопасно)
// И срабатывать на опасные варианты:
// password: "admin123" - реальный пароль (опасно!)
// secret: "mysecret" - реальный секрет (опасно!)
func isEnvReference(value string) bool {
	lower := strings.ToLower(value)
	return strings.Contains(lower, "env") ||
		strings.Contains(lower, "${") ||
		strings.HasPrefix(lower, "$") ||
		strings.Contains(lower, "env") ||
		strings.Contains(lower, "file://") ||
		strings.Contains(lower, "vault")
}

// isSafeKey проверяет, что ключ указывает на безопасный способ хранения
func isSafeKey(key string) bool {
	lower := strings.ToLower(key)

	// Ключи, которые указывают на переменные окружения или файлы
	safePatterns := []string{
		"_env",  // password_env
		"_file", // password_file
		"_path", // password_path
		"_from", // secret_from
	}

	for _, pattern := range safePatterns {
		if strings.HasSuffix(lower, pattern) {
			return true
		}
	}

	return false
}
