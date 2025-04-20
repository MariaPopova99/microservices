package main

import (
	"context"
	"flag"
	"log"
	"net"

	configg "github.com/MariaPopova99/microservices/internal/config"
	config "github.com/MariaPopova99/microservices/internal/config/env"
	"github.com/MariaPopova99/microservices/internal/converter"
	urls "github.com/MariaPopova99/microservices/internal/repository/urls"
	"github.com/MariaPopova99/microservices/internal/service"
	serv "github.com/MariaPopova99/microservices/internal/service/urls"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

//const grpcPort = 50051

type server struct {
	desc.UnimplementedLongShortV1Server
	urlService service.LongShortService
}

func (s *server) GetShort(ctx context.Context, req *desc.GetShortRequest) (*desc.GetShortResponse, error) {
	shortUrl, err := s.urlService.GetShort(ctx, converter.ToLongUrlsFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &desc.GetShortResponse{
		Short_Url: shortUrl.ShortUrl,
		CreatedAt: timestamppb.New(shortUrl.CreatedAt),
	}, nil
}

func (s *server) GetLong(ctx context.Context, req *desc.GetLongRequest) (*desc.GetLongResponse, error) {
	longtUrl, err := s.urlService.GetLong(ctx, converter.ToShortUrlsFromDesc(req))
	if err != nil {
		return nil, err
	}
	return &desc.GetLongResponse{
		Long_Url:  longtUrl.LongUrl,
		CreatedAt: timestamppb.New(longtUrl.CreatedAt),
	}, nil
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
	desc.RegisterLongShortV1Server(s, &server{urlService: ursServ})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
