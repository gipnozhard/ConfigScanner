BINARY = config-scanner

build:
	go build -o $(BINARY) ./cmd/cli

run: build
	./$(BINARY) examples/test.yaml

run-bad: build
	./$(BINARY) examples/bad-config.json

run-good: build
	./$(BINARY) examples/good-config.json

run-stdin: build
	cat $(file) | ./$(BINARY) --stdin

clean:
	rm -f $(BINARY)

run-http:
	go run cmd/httpserver/http.go

run-http-port:
	go run cmd/httpserver/http.go -port=$(port)

test-http:
	@echo "=== Тестирование HTTP сервера ==="
	@echo ""
	@echo "1. Health check:"
	@curl -s http://localhost:8080/health
	@echo ""
	@echo "2. Анализ bad-config.json:"
	@curl -s -X POST http://localhost:8080/analyze -H "Content-Type: application/json" -d @examples/bad-config.json
	@echo ""
	@echo "3. Анализ good-config.json:"
	@curl -s -X POST http://localhost:8080/analyze -H "Content-Type: application/json" -d @examples/good-config.json
	@echo ""
	@echo "4. Анализ test.yaml:"
	@curl -s -X POST http://localhost:8080/analyze --data-binary @examples/test.yaml
	@echo ""

stop-http:
	@echo "Остановка HTTP сервера на порту 8080..."
	@lsof -ti:8080 | xargs kill -9 2>/dev/null || echo "Сервер остановлен!"

run-dir:
	go run cmd/cli/main.go --dir $(DIR)

run-dir-silent:
	go run cmd/cli/main.go --dir --silent $(DIR)

help:
	@echo "Команды cli:"
	@echo "  make build      - собрать программу"
	@echo "  make run        - запустить на test.yaml"
	@echo "  make run-bad    - запустить на bad-config.json"
	@echo "  make run-good   - запустить на good-config.json"
	@echo "  make run-stdin  - запустить из STDIN"
	@echo "	 make run-dir    - рекурсивно анализировать директорию (make run-dir DIR=./configs)"
	@echo "  make clean      - удалить бинарник"

	@echo "Команды http:"
	@echo "  make run-http        - запустить HTTP сервер (порт 8080)"
	@echo "  make run-http-port   - запустить HTTP сервер с кастомным портом"
	@echo "  make build-http      - собрать бинарник HTTP сервера"
	@echo "  make test-http       - протестировать HTTP сервер (отправить тестовые запросы)"
	@echo "  make stop-http       - остановить HTTP сервер (если запущен в фоне)"

