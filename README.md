# **ConfigScanner** 
— утилита командной строки для анализа конфигурационных файлов веб-приложений (YAML/JSON). Автоматически выявляет потенциально опасные настройки безопасности и выдаёт рекомендации по их устранению.

ConfigScanner анализирует конфигурационные файлы и находит:

- **Debug-логирование** — уровень `debug` в настройках логирования
- **Пароли в открытом виде** — хранение паролей, секретов и токенов без шифрования
- **Прослушивание всех интерфейсов** — использование `0.0.0.0` без ограничений
- **Отключённый TLS** — шифрование выключено в production
- **Слабые алгоритмы** — MD5, DES, RC4, SHA-1 и другие устаревшие алгоритмы

Для каждой проблемы выводится:
- Уровень опасности (`LOW`, `MEDIUM`, `HIGH`)
- Краткое объяснение
- Рекомендация по исправлению

## Возможности

- **CLI утилита** - анализ файлов конфигурации
- **HTTP сервер** - REST API для анализа конфигураций
- **gRPC сервер** - высокопроизводительный RPC API для анализа конфигураций
- **Поддержка JSON и YAML** форматов
- **5 правил безопасности**:
    - Debug-логирование (`LOW`)
    - Пароли в открытом виде (`HIGH`)
    - Прослушивание 0.0.0.0 (`MEDIUM`)
    - Отключённый TLS (`HIGH`)
    - Слабые алгоритмы (`HIGH`)

## Установка
```bash
git clone https://github.com/gipnozhard/ConfigScanner.git
cd ConfigScanner
go build -o config-scanner ./cmd/cli
```

## С помощью Makefile
```bash
git clone https://github.com/gipnozhard/ConfigScanner.git
cd ConfigScanner
make build
./bin/config-scanner --help
```

## Запуск CLI-приложения
### через терминал:
```bash
go run cmd/cli/main.go <config-file>
```
#### YAML файл:
```bash
go run cmd/cli/main.go examples/test.yaml
```
#### JSON файл:
#### BAD JSON:
```bash
go run cmd/cli/main.go examples/bad-config.json
```
#### GOOD JSON:
```bash
go run cmd/cli/main.go examples/good-config.json
```

Эти конфигурационные тестовые файлы находятся в детриктории examples.

### Команды Makefile:
 * make build -- Собрать бинарный файл в bin/config-scanner
 * make run	-- Собрать и запустить на test.yaml
 * make run-bad -- Запустить на bad-config.json
 * make run-good -- Запустить на good-config.json
 * make run-stdin <config-file> -- Запустить в режиме STDIN
 * make clean -- Удалить бинарные файлы
 * make help -- Показать справку

#### make help описание:
* make build      - собрать программу
* make run        - запустить на test.yaml
* make run-bad    - запустить на bad-config.json
* make run-good   - запустить на good-config.json
* make run-stdin <config-file>  - запустить из STDIN
* make clean      - удалить бинарник

### через Makefile:
```bash
make run <config-file>
```
#### YAML файл:
```bash
make run examples/test.yaml
```
#### JSON файл:
#### BAD JSON:
```bash
make run examples/bad-config.json
```
#### GOOD JSON:
```bash
make run examples/good-config.json
```

### Пример вывода если есть ошибки и вывод, так же с ошибкой (exit status 1) :
```text
Найдены потенциальные проблемы:

HIGH: plain-password. Обнаружен пароль/секрет в открытом виде. Используйте переменные окружения.
HIGH: weak-algorithm. слишком слабый алгоритм: md5 (MD5 устарел и имеет коллизии). Замените его на более безопасный.
MEDIUM: bind-all. Сервер слушает на всех интерфейсах (0.0.0.0). Ограничьте доступ, используйте localhost или конкретный IP-адрес.
HIGH: tls-disabled. TLS отключен. Включите TLS для безопасного соединения.
LOW: debug-log. Обнаружен debug-режим логирования. Поменяйте режим на более избирательный (info+).
exit status 1
```

### Пример вывода если нет ошибок:
```text
Проблем не найдено! Конфигурация безопасна.
```

### Опции
* -s, --silent -- Не выходить с ошибкой при наличии проблем
* --stdin -- Читать конфигурацию из STDIN вместо файла

### Пример вывода -s, --silent:
#### Напрямую через терминал
```bash
go run cmd/cli/main.go -s <config-file>
go run cmd/cli/main.go --silent <config-file>
```

#### Использовать Makefile
```bash
make run --stdin <config-file>
make run --silent <config-file>
```
### Читать конфигурацию из STDIN:
#### Напрямую через терминал
```bash
go run cmd/cli/main.go --stdin <config-file>
```

#### Использовать Makefile
```bash
make run-stdin file=<config-file>
```

## Структура приложения
```
ConfigScanner
├── cmd
│   ├── cli                   # Точка входа CLI
│   │   └── main.go           # Главный файл с логикой флагов и запуска
│   ├── grpc                  # grpc 
│   │   ├── client
│   │   │   └── client.go
│   │   └── server
│   │       └── server.go
│   └── httpserver            # HTTP сервер
│       └── http.go           # Запуск сервера
├── examples                  # Примеры конфигурационных файлов
│   ├── bad-config.json       # Пример JSON конфига (опасный)
│   ├── good-config.json      # Пример JSON конфига (безопасный)
│   └── test.yaml             # Пример YAML конфига (опасный)
├── go.mod                    # Go модуль
├── go.sum                    # Контрольные суммы зависимостей
├── internal                  # Внутренние пакеты (не для внешнего использования)
│   ├── algo_config           # Рекурсивный анализ директории
│   │   └── cli
│   │       └── cli.go
│   ├── analyzer              # Анализатор конфигурации
│   │   └── analyzer.go       # Запускает все правила на проверку
│   ├── httphandlers
│   │   └── handlers.go       # Обработчики запросов
│   ├── output                # Форматирование вывода
│   │   └── output.go         # TextFormatter, интерфейс Formatter
│   ├── parser                # Парсер конфигураций 
│   │   └── parser.go         # Определяет формат (YAML/JSON) и парсит
│   └── rules                 # Правила проверки
│       ├── bind_all.go       # Правило: прослушка 0.0.0.0 (MEDIUM)
│       ├── debug_log.go      # Правило: debug-логирование (LOW)
│       ├── open_password.go  # Правило: пароли в открытом виде (HIGH)
│       ├── rule              # Интерфейс правила
│       │   └── rule.go       # Определяет контракт Rule
│       ├── tls_disabled.go   # Правило: отключённый TLS (HIGH)
│       └── weak_algorithm.go # Правило: слабые алгоритмы (HIGH)
├── Makefile                  # Автоматизация сборки
├── pkg                       # Публичные пакеты
│   └── models                # Модели данных
│       └── problem.go        # Problem, LevelProblem (LOW/MEDIUM/HIGH)
├── proto                     # протофайл
│   ├── grpc_grpc.pb.go
│   ├── grpc.pb.go
│   └── grpc.proto
└── README.md                 # Документация

13 directories, 19 files
```

## Архитектура приложения
```
┌─────────────┐
│   main.go   │ ← Парсинг флагов, чтение файла/STDIN
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Parser     │ ← Определение формата (YAML/JSON) │
└──────┬──────┘
       │ config (map[string]interface{})
       ▼
┌─────────────┐
│  Analyzer   │ ← Запуск всех правил │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│           Rules (5 шт.)             │
├─────────────────────────────────────┤
│ DebugLogRule    → LOW               │
│ PlainPasswordRule → HIGH            │
│ BindAllRule     → MEDIUM            │
│ TLSDisabledRule → HIGH              │
│ WeakAlgorithmRule → HIGH            │
└──────┬──────────────────────────────┘
       │ []models.Problem
       ▼
┌─────────────┐
│   Output    │ ← Форматирование и вывод  │
└─────────────┘
       │
       ▼
   Console (stdout)
```

## Что именно считается "проблемой"

Утилита реагирует на любое из 5 правил:

| Правило|Условие срабатывания|Пример
|--------|---------|----------|
|debug-log|level: debug|logging.level = "debug"
|plain-password|Пароль/секрет в открытом виде|password: "admin123"
|bind-all	|host: 0.0.0.0	 |server.host = "0.0.0.0"
|tls-disabled	|tls.enabled: false	 |tls.enabled = false
|weak-algorithm	|MD5, DES, RC4 и др.	 |hash_algorithm: "md5"


# Запускать утилиту в качестве http-сервера, который имеет REST API. 
## С помощью API можно также передавать конфигурацию на проверку и получать в качестве ответа массив проблем.

### Запуск

```bach
# Стандартный запуск (порт 8080)
make run-http

# Или напрямую
go run cmd/httpserver/http.go

# С кастомным портом
go run cmd/httpserver/http.go -port=9090
```

### API Эндпоинты
#### GET /health - проверка здоровья сервера
Пример запроса:
```bach
curl http://localhost:8080/health
```
Пример ответа:
```json
{
  "status": "ok"
}
```
#### POST /analyze - анализ конфигурации
Пример запроса с JSON файлом:
```bach
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d @examples/bad-config.json
```
Пример запроса с YAML файлом:
```bach
curl -X POST http://localhost:8080/analyze \
  --data-binary @examples/test.yaml
```
Пример запроса с прямой JSON строкой:
```bach
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{"logging":{"level":"debug"},"server":{"host":"0.0.0.0"}}'
```

Пример ответа (есть проблемы):
```json
{
  "count": 2,
  "problems": [
    {
      "path": "logging.level",
      "explanation": "Обнаружен debug-режим логирования",
      "rule": "debug-log",
      "level_problem": "LOW",
      "recommendation": "Поменяйте режим на более избирательный (info+)"
    },
    {
      "path": "server.host",
      "explanation": "Сервер слушает на всех интерфейсах (0.0.0.0)",
      "rule": "bind-all",
      "level_problem": "MEDIUM",
      "recommendation": "Ограничьте доступ, используйте localhost или конкретный IP-адрес"
    }
  ]
}
```

Пример ответа (проблем нет):
```json
{
  "count": 0,
  "problems": null
}
```

## Коды ответов HTTP
|Код	| Описание                              
|-|-|
|200	| Успешная обработка (проблемы могут быть или отсутствовать) 
|400	|Ошибка парсинга конфига (невалидный JSON/YAML)
|405	|Использован неподдерживаемый метод (нужен POST)

## Тестирование HTTP сервера

```bash
# Запустить тесты (сервер должен быть запущен)
make test-http
```

## Команды Makefile для http

|Команда	|Описание
|- |-
|make build-http	|Собрать HTTP сервер
|make run-http	|Запустить HTTP сервер (порт 8080)
|make run-http-port |port=9090	Запустить HTTP сервер на кастомном порту
|make test-http	|Протестировать HTTP сервер
|make stop-http	|Остановить HTTP сервер
|make clean	|Удалить бинарные файлы


# Рекурсивный анализ директории

Утилита может рекурсивно анализировать все конфигурационные файлы в директории и её поддиректориях.

1. Проверяет, существует ли указанная директория

2. Рекурсивно обходит все вложенные папки

3. Находит файлы с расширениями .json, .yaml, .yml

4. Анализирует каждый найденный файл

5. Собирает статистику по всем файлам

### Использование
```bash
# Анализ директории
go run cmd/cli/main.go --dir ./configs

# С silent режимом (не выходить с ошибкой)
go run cmd/cli/main.go --dir --silent ./configs

# Через Makefile
make run-dir DIR=./examples
```

### Пример вывода
```bash

Рекурсивный анализ директории: ./examples
Найдено конфигурационных файлов: 3

════════════════════════════════════════════════════════════

examples/bad-config.json
────────────────────────────────────────
Найдены потенциальные проблемы:

LOW: debug-log. Обнаружен debug-режим логирования. Поменяйте режим на более избирательный (info+).
HIGH: plain-password. Обнаружен пароль/секрет в открытом виде. Используйте переменные окружения.
MEDIUM: bind-all. Сервер слушает на всех интерфейсах (0.0.0.0). Ограничьте доступ, используйте localhost или конкретный IP-адрес.
HIGH: tls-disabled. TLS отключен. Включите TLS для безопасного соединения.
HIGH: weak-algorithm. слишком слабый алгоритм: des (DES слишком слабый). Замените его на более безопасный.
HIGH: weak-algorithm. слишком слабый алгоритм: md5 (MD5 устарел и имеет коллизии). Замените его на более безопасный.

examples/test.yaml
────────────────────────────────────────
Найдены потенциальные проблемы:

LOW: debug-log. Обнаружен debug-режим логирования. Поменяйте режим на более избирательный (info+).
HIGH: plain-password. Обнаружен пароль/секрет в открытом виде. Используйте переменные окружения.
HIGH: plain-password. Обнаружен пароль/секрет в открытом виде. Используйте переменные окружения.
MEDIUM: bind-all. Сервер слушает на всех интерфейсах (0.0.0.0). Ограничьте доступ, используйте localhost или конкретный IP-адрес.
HIGH: weak-algorithm. слишком слабый алгоритм: md5 (MD5 устарел и имеет коллизии). Замените его на более безопасный.
HIGH: weak-algorithm. слишком слабый алгоритм: rc4 (RC4 небезопасен). Замените его на более безопасный.
════════════════════════════════════════════════════════════

Статистика:
    Всего файлов: 3
    Файлов с проблемами: 2
    Всего проблем: 12
exit status 1
```

## Структура вывода:

|Часть вывода	|Описание
|- |-
|Заголовок	|Путь к директории и количество найденных файлов
|Разделитель	|Визуальное отделение (символ ═)
|Отчёт по файлу	|Для каждого файла с проблемами: имя файла и список проблем
|Статистика	|Общее количество файлов, файлов с проблемами, всего проблем

# gRPC сервер

Высокопроизводительный gRPC API для анализа конфигураций. Использует protobuf для сериализации данных.

## Запуск
```bash
# Стандартный запуск (порт 50051)
make run-grpc-server

# С кастомным портом
make run-grpc-server-port port=9090

# Или напрямую
go run cmd/grpc/server/server.go
go run cmd/grpc/server/server.go -port=50051
```
## gRPC Клиент
```bash
# Запуск клиента для анализа файла
make run-grpc-client file=examples/bad-config.json

# С указанием адреса сервера
go run cmd/grpc/client/client.go -server=localhost:50051 -file=examples/bad-config.json
```

## Пример вывода gRPC клиента
```bash
$ make run-grpc-client file=examples/bad-config.json

Результат анализа: examples/bad-config.json
═══════════════════════════════════════════════════════════════
Найдено проблем: 6

1. [LOW] debug-log
   Путь: logging.level
   Проблема: Обнаружен debug-режим логирования
   Рекомендация: Поменяйте режим на более избирательный (info+)

2. [HIGH] plain-password
   Путь: database.password
   Проблема: Обнаружен пароль/секрет в открытом виде
   Рекомендация: Используйте переменные окружения

3. [MEDIUM] bind-all
   Путь: server.host
   Проблема: Сервер слушает на всех интерфейсах (0.0.0.0)
   Рекомендация: Ограничьте доступ, используйте localhost или конкретный IP-адрес

4. [HIGH] tls-disabled
   Путь: server.tls.enabled
   Проблема: TLS отключен
   Рекомендация: Включите TLS для безопасного соединения

5. [HIGH] weak-algorithm
   Путь: security.hash_algorithm
   Проблема: слишком слабый алгоритм: md5 (MD5 устарел и имеет коллизии)
   Рекомендация: Замените его на более безопасный

6. [HIGH] weak-algorithm
   Путь: security.encryption
   Проблема: слишком слабый алгоритм: des (DES слишком слабый)
   Рекомендация: Замените его на более безопасный

═══════════════════════════════════════════════════════════════
```

## Команды Makefile для gRPC

|Команда	|Описание
|- |-
|make run-grpc-server	|Запустить gRPC сервер (порт 50051)
|make run-grpc-server-port port=9090	|Запустить gRPC сервер на кастомном порту
|make run-grpc-client file=./config.json	|Запустить gRPC клиент для анализа файла
|make build-grpc	|Собрать бинарники gRPC сервера и клиента
|make proto-gen	|Сгенерировать protobuf код из .proto файла