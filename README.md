**ConfigScanner** — утилита командной строки для анализа конфигурационных файлов веб-приложений (YAML/JSON). Автоматически выявляет потенциально опасные настройки безопасности и выдаёт рекомендации по их устранению.

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
│   └── cli                   # Точка входа CLI
│       └── main.go           # Главный файл с логикой флагов и запуска
├── examples                  # Примеры конфигурационных файлов
│   ├── bad-config.json       # Пример JSON конфига (опасный)
│   ├── good-config.json      # Пример JSON конфига (безопасный)
│   └── test.yaml             # Пример YAML конфига (опасный)
├── go.mod                    # Go модуль
├── go.sum                    # Контрольные суммы зависимостей
├── internal                  # Внутренние пакеты (не для внешнего использования)
│   ├── analyzer              # Анализатор конфигурации
│   │   └── analyzer.go       # Запускает все правила на проверку
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