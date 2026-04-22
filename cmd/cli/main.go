package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"ConfigScanner/internal/analyzer"
	"ConfigScanner/internal/output"
	"ConfigScanner/internal/parser"
)

func main() {
	// Парсим флаги
	// flag.Bool() создаёт флаг типа bool (true/false)
	// Параметры:
	//   1. "s" - короткое имя флага (один дефис: -s)
	//   2. false - значение по умолчанию (флаг по умолчанию выключен)
	//   3. "Не выходить..." - описание для справки (--help)
	// Возвращает: *bool (указатель на bool)
	silent := flag.Bool("s", false, "Не выходить с ошибкой при наличии проблем")

	// То же самое, но длинная версия флага (--silent)
	// Это удобно: можно писать и -s, и --silent
	silentLong := flag.Bool("silent", false, "Не выходить с ошибкой при наличии проблем")

	// Флаг для чтения из STDIN вместо файла
	stdin := flag.Bool("stdin", false, "Читать конфигурацию из STDIN вместо файла")

	// flag.Parse() - ОБЯЗАТЕЛЬНЫЙ вызов!
	// Он анализирует os.Args (все аргументы командной строки) и заполняет созданные флаги
	// БЕЗ ЭТОГО флаги всегда будут иметь значение по умолчанию (false)
	flag.Parse()

	// Объединяем короткую и длинную версии silent-флага
	// Если указан хотя бы один из них, isSilent = true
	// *silent - разыменование указателя (получаем значение bool)
	isSilent := *silent || *silentLong

	// Читаем конфигурацию
	var configData []byte
	var err error

	// Проверяем: указан ли флаг --stdin?
	if *stdin {
		// Режим чтения из STDIN (стандартный поток ввода)
		// Используется в цепочках: cat config.json | программа --stdin
		// io.ReadAll читает всё из os.Stdin (стандартный ввод)
		// os.Stdin - это файловый дескриптор, связанный с вводом с клавиатуры
		configData, err = io.ReadAll(os.Stdin)
		if err != nil {
			// Если ошибка - пишем в STDERR (стандартный поток ошибок)
			// %v - формат для вывода значения ошибки
			fmt.Fprintf(os.Stderr, "Ошибка чтения из STDIN: %v\n", err)
			os.Exit(1) // Выход с кодом 1 (ошибка)
		}
	} else {
		// flag.NArg() возвращает количество позиционных аргументов (не флагов)
		if flag.NArg() < 1 { // если число аргументов NArg меньше одного, то выводим ошибку
			fmt.Fprintln(os.Stderr, "укажите файл конфигурации")
			os.Exit(1)
		}

		// flag.Arg(0) - берём первый позиционный аргумент (индекс 0)
		// Это имя файла, который нужно прочитать
		filename := flag.Arg(0)
		// os.ReadFile читает весь файл целиком в байтовый срез
		configData, err = os.ReadFile(filename)
		if err != nil {
			// Не удалось прочитать файл (нет прав, файл не существует и т.д.)
			fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filename, err)
			os.Exit(1)
		}
	}

	// Парсим конфиг
	// Создаём новый парсер (из нашего пакета parser)
	// NewParser() - конструктор, возвращает готовый объект парсера
	configParser := parser.NewParser()
	// Парсим байтовые данные в структуру map[string]interface{}
	// Парсер сам определяет формат: JSON или YAML
	// Возвращает: config: map с распарсенными данными, err: ошибка, если формат неправильный
	config, err := configParser.Parse(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга: %v\n", err)
		os.Exit(1)
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
	// Выводим результат в консоль
	// Print внутри вызывает formatter.Format() и печатает в STDOUT
	output.Print(problems, formatter)

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
