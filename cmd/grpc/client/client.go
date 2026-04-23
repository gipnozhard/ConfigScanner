package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"ConfigScanner/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Алгоритм работы:
//  1. Парсит флаги командной строки
//  2. Читает файл с конфигурацией
//  3. Подключается к gRPC серверу
//  4. Отправляет запрос на анализ
//  5. Выводит полученный результат
func main() {
	// Флаг для указания адреса gRPC сервера
	// По умолчанию: localhost:50051 (стандартный порт gRPC)
	serverAddr := flag.String("server", "localhost:50051", "gRPC сервер адрес")

	// Флаг для указания файла с конфигурацией
	// Обязательный флаг (проверяем после парсинга)
	filename := flag.String("file", "", "файл для анализа")

	// flag.Parse() - ОБЯЗАТЕЛЬНЫЙ вызов!
	// Анализирует os.Args и заполняет созданные флаги
	// Без этого флаги всегда будут иметь значения по умолчанию
	flag.Parse()

	// Проверяем, указан ли файл
	// Если filename всё ещё пустая строка - значит пользователь забыл указать -file
	if *filename == "" {
		// log.Fatal - выводит сообщение и завершает программу с кодом 1
		log.Fatal("укажите файл: --file config.json")
	}

	// Читаем файл
	data, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("Ошибка чтения: %v", err)
	}

	// Подключаемся к gRPC серверу
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	// defer - закрываем соединение при выходе из функции
	// Важно: всегда закрывать соединение, чтобы не было утечек
	defer conn.Close()

	// Создаём клиента для вызова RPC методов
	// proto.NewConfigAnalyzerClient - функция, сгенерированная protoc
	// Принимает соединение, возвращает объект клиента
	client := proto.NewConfigAnalyzerClient(conn)

	// Отправляем запрос
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := client.Analyze(ctx, &proto.AnalyzeRequest{
		ConfigData: data,
		Filename:   *filename,
	})
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Выводим заголовок с именем файла
	fmt.Printf("\nРезультат анализа: %s\n", *filename)
	// strings.Repeat("═", 50) - но мы используем прямую строку
	// Символ ═ (двойная линия) используется для рамки
	fmt.Println("═══════════════════════════════════════════════════════════════")

	// Проверяем количество проблем
	if resp.Count == 0 {
		fmt.Println("Проблем не найдено!")
	} else {
		// Если проблемы есть - выводим каждую
		fmt.Printf("Найдено проблем: %d\n\n", resp.Count)

		// Проходим по всем проблемам из ответа
		// resp.Problems - это срез указателей на Problem
		for i, p := range resp.Problems {
			fmt.Printf("%d. [%s] %s\n", i+1, p.LevelProblem, p.Rule) // Номер проблемы, уровень опасности и название правила
			fmt.Printf("   Путь: %s\n", p.Path)                      // Строка 2: Путь к проблемному полю в конфиге
			fmt.Printf("   Проблема: %s\n", p.Explanation)           // Строка 3: Описание проблемы
			fmt.Printf("   Рекомендация: %s\n\n", p.Recommendation)  // Строка 4: Рекомендация по исправлению
		}
	}
	fmt.Println("═══════════════════════════════════════════════════════════════")
}
