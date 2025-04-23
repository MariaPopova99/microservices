package main

import (
	"context"
	"flag"
	"log"
	"net"

	urlsApi "github.com/MariaPopova99/microservices/internal/api/urls"
	configg "github.com/MariaPopova99/microservices/internal/config"
	config "github.com/MariaPopova99/microservices/internal/config/env"
	urls "github.com/MariaPopova99/microservices/internal/repository/urls"
	serv "github.com/MariaPopova99/microservices/internal/service/urls"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := configg.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	urlRepo := urls.NewRepository(pool)
	ursServ := serv.NewService(urlRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterLongShortV1Server(s, urlsApi.NewImplementation(ursServ))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
