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

help:
	@echo "Команды:"
	@echo "  make build      - собрать программу"
	@echo "  make run        - запустить на test.yaml"
	@echo "  make run-bad    - запустить на bad-config.json"
	@echo "  make run-good   - запустить на good-config.json"
	@echo "  make run-stdin  - запустить из STDIN"
	@echo "  make clean      - удалить бинарник"