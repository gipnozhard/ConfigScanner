package httphandlers

import (
	"encoding/json"
	"io"
	"net/http"

	"ConfigScanner/internal/analyzer"
	"ConfigScanner/internal/parser"
)

// HealthHandler обрабатывает GET /health запросы
// проверка работоспособности сервера (health check)
// Ответ: всегда возвращает 200 OK и JSON {"status":"ok"}
// Даже если сервер перегружен, этот эндпоинт должен отвечать быстро
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type ответ будет в формате JSON
	// Без этого можно не понять, как парсить ответ
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем HTTP статус 200 OK
	// 200 означает "всё хорошо, запрос обработан успешно"
	w.WriteHeader(http.StatusOK)

	// Отправляем тело ответа
	w.Write([]byte(`{"status":"ok"}`))
}

// AnalyzeHandler обрабатывает POST /analyze запросы
// Что делает:
// 1. Проверяет HTTP метод (должен быть POST)
// 2. Читает тело запроса (конфиг в формате JSON или YAML)
// 3. Парсит конфиг во внутреннее представление (map)
// 4. Запускает анализ по всем правилам безопасности
// 5. Возвращает список проблем в формате JSON
// Пример запроса:
//
//	curl -X POST http://localhost:8080/analyze \
//	  -H "Content-Type: application/json" \
//	  -d '{"logging":{"level":"debug"}}'
//
// Пример ответа:
//
//	{
//	  "problems": [
//	    {
//	      "rule": "debug-log",
//	      "level_problem": "LOW",
//	      "explanation": "Обнаружен debug-режим логирования",
//	      "recommendation": "Поменяйте режим на более избирательный (info+)"
//	    }
//	  ],
//	  "count": 1
//	}
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Только POST
	if r.Method != http.MethodPost {
		http.Error(w, "только POST метод", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	// io.ReadAll - читает всё содержимое до конца, body: срез байтов ([]byte) с данными
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ошибка чтения тела: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Парсим конфиг (JSON или YAML)
	configParser := parser.NewParser()
	config, err := configParser.Parse(body)
	if err != nil {
		http.Error(w, "ошибка парсинга: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Анализируем
	analyzer := analyzer.NewAnalyzer()
	problems := analyzer.Analyze(config)

	// Отдаем результат
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"problems": problems,
		"count":    len(problems),
	})
}
