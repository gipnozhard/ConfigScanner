package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"ConfigScanner/internal/analyzer"
	"ConfigScanner/internal/parser"
	"ConfigScanner/proto"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedConfigAnalyzerServer
}

// Analyze - реализация gRPC метода
func (s *server) Analyze(ctx context.Context, req *proto.AnalyzeRequest) (*proto.AnalyzeResponse, error) {
	log.Printf("Получен запрос от %s", req.Filename)

	// Парсим конфиг (переиспользуем существующий парсер)
	configParser := parser.NewParser()
	config, err := configParser.Parse(req.ConfigData)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга: %v", err)
	}

	// Анализируем (переиспользуем существующий анализатор)
	analyzer := analyzer.NewAnalyzer()
	problems := analyzer.Analyze(config)

	// Конвертируем проблемы в protobuf формат
	protoProblems := make([]*proto.Problem, len(problems))
	for i, p := range problems {
		protoProblems[i] = &proto.Problem{
			Path:           p.Path,
			Explanation:    p.ParseV,
			Rule:           p.Rule,
			LevelProblem:   string(p.LevelProblem),
			Recommendation: p.Recommendation,
		}
	}

	return &proto.AnalyzeResponse{
		Problems: protoProblems,
		Count:    int32(len(problems)),
	}, nil
}

func main() {
	port := flag.Int("port", 50051, "gRPC сервер порт")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Ошибка запуска: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterConfigAnalyzerServer(s, &server{})

	log.Printf("gRPC сервер запущен на порту %d", *port)
	log.Printf("Методы:")
	log.Printf("  Analyze - анализ конфигурации")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}
