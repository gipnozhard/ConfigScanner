package rules

import (
	"ConfigScanner/pkg/models"
	"strings"
)

// WeakAlgorithmRule структура для правила "слабый алгоритм шифрования/хеширования"
// Это правило ищет устаревшие и небезопасные криптографические алгоритмы
//
// Чем это опасно:
// Устаревшие алгоритмы имеют известные уязвимости
// Взломщики могут подобрать коллизии (MD5, SHA-1)
// Слишком короткие ключи можно перебрать (DES - 56 бит)
// Примеры опасных конфигов:
//
//	security:
//	  password_hash: md5      # ← MD5 легко взломать
//	  encryption: des         # ← DES перебирается за часы
//	  cipher: rc4            # ← RC4 имеет уязвимости
//
// Какие алгоритмы считаются слабыми:
//   - Хеши: MD5, SHA-1, MD4, MD2
//   - Шифрование: DES, 3DES, RC4, Blowfish (старый)
//   - Другие: все что короче 128 бит
type WeakAlgorithmRule struct{}

func (r WeakAlgorithmRule) Name() string {
	return "weak-algorithm"
}

func (r WeakAlgorithmRule) Check(config map[string]interface{}, path string) []models.Problem {
	var problems []models.Problem

	// Словарь слабых алгоритмов с пояснениями, почему они плохие
	// ключ: название алгоритма в нижнем регистре
	// значение: причина, почему алгоритм небезопасен
	weakAlgorithms := map[string]string{
		"md5":      "MD5 устарел и имеет коллизии",
		"sha1":     "SHA-1 устарел и небезопасен",
		"des":      "DES слишком слабый",
		"rc4":      "RC4 небезопасен",
		"blowfish": "Blowfish устарел",
		"md4":      "MD4 полностью сломан",
		"md2":      "MD2 слишком старый",
		"3des":     "3DES медленный и устаревает",
		"sha0":     "SHA-0 никогда не был безопасным",
		"ripemd":   "RIPEMD устарел",
	}

	for key, value := range config {
		currentPath := joinPath(path, key)

		// Проверяем ключи, связанные с алгоритмами
		if isAlgorithmKey(key) {
			if str, ok := value.(string); ok {
				lowerStr := strings.ToLower(str)
				for weakAlgo, reason := range weakAlgorithms {
					if strings.Contains(lowerStr, weakAlgo) {
						problems = append(problems, models.Problem{
							Rule:           r.Name(),
							LevelProblem:   models.High, // Использование слабых алгоритмов может привести к компрометации данных!
							Path:           currentPath,
							ParseV:         "слишком слабый алгоритм: " + weakAlgo + " (" + reason + ")",
							Recommendation: "Замените его на более безопасный",
						})
						break
					}
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

// isAlgorithmKey проверяет, является ли ключ связанным с криптографическими алгоритмами
func isAlgorithmKey(key string) bool {
	lower := strings.ToLower(key)
	// Проверяем наличие ключевых слов, указывающих на алгоритмы
	return strings.Contains(lower, "algorithm") ||
		strings.Contains(lower, "hash") ||
		strings.Contains(lower, "cipher") ||
		strings.Contains(lower, "encryption") ||
		strings.Contains(lower, "digest")
}
