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
	silent := flag.Bool("s", false, "Не выходить с ошибкой при наличии проблем")
	silentLong := flag.Bool("silent", false, "Не выходить с ошибкой при наличии проблем")
	stdin := flag.Bool("stdin", false, "Читать конфигурацию из STDIN вместо файла")
	flag.Parse()

	isSilent := *silent || *silentLong

	// Читаем конфигурацию
	var configData []byte
	var err error

	if *stdin {
		configData, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения из STDIN: %v\n", err)
			os.Exit(1)
		}
	} else {
		if flag.NArg() < 1 {
			printUsage()
			os.Exit(1)
		}

		filename := flag.Arg(0)
		configData, err = os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filename, err)
			os.Exit(1)
		}
	}

	// Парсим конфиг
	configParser := parser.NewParser()
	config, err := configParser.Parse(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга: %v\n", err)
		os.Exit(1)
	}

	// Анализируем
	analyzer := analyzer.NewAnalyzer()
	problems := analyzer.Analyze(config)

	// Выводим результат
	formatter := output.TextFormatter{}
	output.Print(problems, formatter)

	// Выходим с кодом ошибки при наличии проблем
	if len(problems) > 0 && !isSilent {
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Использование: %s [options] <config-file>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Или: cat config | %s --stdin\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nОпции:\n")
	fmt.Fprintf(os.Stderr, "  -s, --silent    Не выходить с ошибкой при наличии проблем\n")
	fmt.Fprintf(os.Stderr, "  --stdin         Читать конфигурацию из STDIN\n")
}
