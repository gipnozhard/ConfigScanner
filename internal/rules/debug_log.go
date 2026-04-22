package rules

import (
	"ConfigScanner/pkg/models"
	"strings"
)

// DebugLogRule структура для правила "debug-логирование"
type DebugLogRule struct{}

// Name возвращает уникальный идентификатор правила
// Метод реализует часть интерфейса Rule
func (r DebugLogRule) Name() string {
	return "debug-log"
}

// Check - основной метод, который проверяет конфигурацию на наличие debug-логирования
//
// Параметры:
//   - config: распарсенный конфиг в виде map (ключ -> значение)
//   - path: текущий путь
//
// Возвращает:
//   - []models.Problem: срез найденных проблем (может быть пустым)
//
// Метод рекурсивно обходит всю структуру конфига в поисках опасных значений
func (r DebugLogRule) Check(config map[string]interface{}, path string) []models.Problem {
	// Создаем пустой срез для хранения найденных проблем
	// Используем var вместо make, так как срез может остаться пустым
	var problems []models.Problem

	// Обход всех ключей текущего уровня конфига
	for key, value := range config {
		// Формируем полный путь до текущего ключа
		// Пример: если path="server" и key="host", то currentPath="server.host"
		// Если path пустой (корневой уровень), то currentPath = key
		currentPath := joinPath(path, key)

		// Проверяем, является ли текущий ключ ключом уровня логирования
		// isLogLevelKey проверяет: "level", "log_level", "loglevel"
		if isLogLevelKey(key) {
			// Пытаемся привести значение к строке (type assertion)
			// str, ok := value.(string) означает:
			//   - str: значение, если оно действительно строка
			//   - ok: true если value это string, false если другой тип
			if str, ok := value.(string); ok && strings.ToLower(str) == "debug" {
				// Условие сработало, значит:
				// 1. value можно привести к строке (ok == true)
				// 2. После приведения к нижнему регистру строка равна "debug"
				//
				// Приводим к нижнему регистру, чтобы находить:
				// "debug", "DEBUG", "Debug", "dEbUg" - все варианты

				// Создаем новую проблему и добавляем в срез
				problems = append(problems, models.Problem{
					Rule:           r.Name(),                            // Имя правила: "debug-log"
					LevelProblem:   models.Low,                          // Уровень опасности: LOW
					Path:           currentPath,                         // Где найдена проблема
					ParseV:         "Обнаружен debug-режим логирования", // Описание
					Recommendation: "Поменяйте режим на более избирательный (info+)",
				})
			}
		}

		// Рекурсивная проверка вложенных объектов
		// Проверяем, является ли значение вложенным в map объектом
		// Пример: value = {"level": "debug"} - это map
		if nested, ok := value.(map[string]interface{}); ok {
			// Если да, вызываем Check рекурсивно для вложенного map
			problems = append(problems, r.Check(nested, currentPath)...)
		}
	}

	// Возвращаем все найденные проблемы
	// Если проблем нет - возвращаем пустой срез (не nil)
	return problems
}

// isLogLevelKey проверяет, является ли ключ ключом уровня логирования
// Функция регистронезависима - "Level", "LEVEL", "level" обрабатываются одинаково
func isLogLevelKey(key string) bool {
	lower := strings.ToLower(key)
	return lower == "level" || lower == "log_level" || lower == "loglevel"
}

// joinPath объединяет путь и ключ в один строковый путь через точку
func joinPath(path, key string) string {
	if path == "" {
		return key
	}
	return path + "." + key
}
