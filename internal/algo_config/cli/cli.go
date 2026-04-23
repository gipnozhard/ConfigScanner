package cli

import (
	"ConfigScanner/internal/analyzer"
	"ConfigScanner/internal/output"
	"ConfigScanner/internal/parser"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// analyzeConfig анализирует один конфиг
// Параметры:
// - configData: байтовый срез с содержимым конфига (JSON или YAML)
// - filename: имя файла (для вывода, может быть пустым при чтении из STDIN)
// - isSilent: если true - не выходить с ошибкой при наличии проблем
// Функция самостоятельно обрабатывает ошибки и вызывает os.Exit() при критических проблемах
func AnalyzeConfig(configData []byte, filename string, isSilent bool) {
	// Парсим конфиг
	// Создаём новый парсер (из нашего пакета parser)
	// NewParser() - конструктор, возвращает готовый объект парсера
	configParser := parser.NewParser()

	// Парсим байтовые данные в структуру map[string]interface{}
	// Парсер сам определяет формат: JSON или YAML
	// Возвращает: config: map с распарсенными данными, err: ошибка, если формат неправильный
	config, err := configParser.Parse(configData)
	if err != nil {
		// Если парсинг не удался - выводим ошибку в STDERR и завершаем программу
		// os.Stderr - стандартный поток ошибок (отдельно от обычного вывода)
		fmt.Fprintf(os.Stderr, "Ошибка парсинга %s: %v\n", filename, err)
		os.Exit(1) // Код 1 означает ошибку
	}

	// Анализируем
	// Создаём анализатор со всеми правилами
	// NewAnalyzer() создаёт анализатор с правилами: DebugLogRule (debug-логирование)
	// PlainPasswordRule (пароли в открытом виде)
	// BindAllRule (0.0.0.0)
	// TLSDisabledRule (отключённый TLS)
	// WeakAlgorithmRule (слабые алгоритмы)
	analyzer := analyzer.NewAnalyzer()

	// Запускаем анализ
	// Analyze рекурсивно обходит весь конфиг и применяет все правила
	// Возвращает срез найденных проблем ([]models.Problem)
	problems := analyzer.Analyze(config)

	// Выводим результат
	// Создаём форматтер для текстового вывода
	// TextFormatter преобразует срез проблем в красивую строку
	formatter := output.TextFormatter{}

	// Если имя файла указано (не STDIN), выводим заголовок
	if filename != "" {
		fmt.Printf("\nАнализ файла: %s\n", filename)
		// strings.Repeat повторяет символ 50 раз → линия-разделитель
		fmt.Println(strings.Repeat("─", 50))
	}

	// Выводим результат в консоль
	// Print внутри вызывает formatter.Format() и печатает в STDOUT
	output.Print(problems, formatter)

	// Добавляем пустую строку после результата (для красоты)
	if filename != "" {
		fmt.Println()
	}

	// Выходим с кодом ошибки при наличии проблем
	// Проверяем:
	// 1. len(problems) > 0 - есть хотя бы одна проблема?
	// 2. !isSilent - не включен ли silent-режим?
	// Если оба условия true - выходим с кодом ошибки 1
	if len(problems) > 0 && !isSilent {
		os.Exit(1)
	}
	// Если проблем нет ИЛИ включен silent-режим os.Exit(0) вызывается неявно
}

// НОВАЯ ФУНКЦИЯ для рекурсивного анализа директории
// analyzeDirectory рекурсивно анализирует все конфиги в директории
// Параметры:
// - dirPath: путь к директории для анализа
// - isSilent: если true - не выходить с ошибкой при наличии проблем
// Функция:
// 1. Проверяет существование директории
// 2. Рекурсивно находит все .json, .yaml, .yml файлы
// 3. Анализирует каждый файл
// 4. Выводит статистику по всем файлам
func AnalyzeDirectory(dirPath string, isSilent bool) {
	// Проверяем, существует ли директория
	// os.Stat возвращает информацию о файле/директории
	info, err := os.Stat(dirPath)
	if err != nil {
		// Ошибка: директория не существует или нет прав доступа
		fmt.Fprintf(os.Stderr, "Ошибка доступа к директории %s: %v\n", dirPath, err)
		os.Exit(1)
	}
	// IsDir() возвращает true, если это директория, а не файл
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "%s не является директорией\n", dirPath)
		os.Exit(1)
	}

	// Собираем все конфиги рекурсивно
	var configFiles []string

	// filepath.Walk - рекурсивный обход всех файлов и папок
	// Параметры:
	// 1. dirPath - откуда начинать обход
	// 2. функция, которая вызывается для каждого элемента
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // передаём ошибку дальше
		}

		// Пропускаем директории, обрабатываем только файлы
		if !info.IsDir() {
			// Получаем расширение файла в нижнем регистре
			// filepath.Ext возвращает расширение (например ".json")
			ext := strings.ToLower(filepath.Ext(path))

			// Проверяем, является ли расширение поддерживаемым
			if ext == ".json" || ext == ".yaml" || ext == ".yml" {
				configFiles = append(configFiles, path) // добавляем в список
			}
		}
		return nil // nil = ошибок нет, продолжаем обход
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка обхода директории: %v\n", err)
		os.Exit(1)
	}

	// Если файлы не найдены - выводим сообщение и выходим
	if len(configFiles) == 0 {
		fmt.Printf("В директории %s не найдено .json, .yaml или .yml файлов\n", dirPath)
		return
	}

	// Заголовок вывода
	fmt.Printf("\nРекурсивный анализ директории: %s\n", dirPath)
	fmt.Printf("Найдено конфигурационных файлов: %d\n\n", len(configFiles))
	// strings.Repeat("═", 60) создаёт линию из 60 символов "═"
	fmt.Println(strings.Repeat("═", 60))

	// Переменные для статистики
	totalProblems := 0 // общее количество проблем во всех файлах
	problemFiles := 0  // количество файлов, в которых есть проблемы

	// Анализируем каждый найденный файл
	for _, file := range configFiles {
		// Читаем содержимое файла
		configData, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Ошибка чтения %s: %v\n", file, err)
			continue // переходим к следующему файлу
		}

		// Парсим конфиг
		configParser := parser.NewParser()
		config, err := configParser.Parse(configData)
		if err != nil {
			fmt.Printf("Ошибка парсинга %s: %v\n", file, err)
			continue // переходим к следующему файлу
		}

		// Анализируем
		analyzer := analyzer.NewAnalyzer()
		problems := analyzer.Analyze(config)

		// Если есть проблемы - выводим их
		if len(problems) > 0 {
			problemFiles++                       // увеличиваем счётчик файлов с проблем
			totalProblems += len(problems)       // добавляем количество проблем
			fmt.Printf("\n%s\n", file)           // выводим имя файла
			fmt.Println(strings.Repeat("─", 40)) // линия-разделитель

			formatter := output.TextFormatter{}
			output.Print(problems, formatter) // выводим проблемы
		}
	}

	// Статистика
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("\nСтатистика:\n")
	fmt.Printf("    Всего файлов: %d\n", len(configFiles))
	fmt.Printf("    Файлов с проблемами: %d\n", problemFiles)
	fmt.Printf("    Всего проблем: %d\n", totalProblems)

	// Если есть проблемы и не включён silent-режим - завершаем с ошибкой
	if totalProblems > 0 && !isSilent {
		os.Exit(1)
	}
}
