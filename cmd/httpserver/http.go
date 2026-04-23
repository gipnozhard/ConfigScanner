package main

import (
	"ConfigScanner/internal/httphandlers"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// flag.String создаёт строковый флаг командной строки
	// Параметры:
	//   1. "port" - имя флага (будет доступен как -port или --port)
	//   2. "8080" - значение по умолчанию, если флаг не указан
	//   3. "порт для сервера" - описание, которое увидит пользователь при --help
	// Возвращает: *string (указатель на строку)
	//
	// Пример использования: go run main.go -port 9090
	port := flag.String("port", "8080", "порт для сервера")

	// flag.Parse() - ОБЯЗАТЕЛЬНЫЙ вызов!
	// Он анализирует os.Args (все аргументы командной строки) и заполняет созданные флаги
	// БЕЗ ЭТОГО ВЫЗОВА флаг port всегда будет иметь значение по умолчанию "8080"
	// Даже если пользователь указал -port 9090
	flag.Parse()

	// Регистрируем обработчики
	// health check эндпоинт (проверка здоровья сервера)
	http.HandleFunc("/health", httphandlers.HealthHandler)

	// Сюда клиенты будут отправлять POST запросы с конфигом
	http.HandleFunc("/analyze", httphandlers.AnalyzeHandler)

	adr := fmt.Sprintf(":%s", *port)
	log.Printf("Сервер запущен на http://localhost%s", adr)
	log.Printf("Доступные эндпоинты:")
	log.Printf("  GET  /health  - проверка сервера пдключения")
	log.Printf("  POST /analyze - анализ конфигурации")

	log.Fatal(http.ListenAndServe(adr, nil))
}
