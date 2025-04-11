package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedLongShortV1Server
}

func (s *server) GetShort(ctx context.Context, req *desc.GetShortRequest) (*desc.GetShortResponse, error) {

	return &desc.GetShortResponse{
		Short_Url: "Short" + req.GetLong_Url(), // Или что-то другое
		CreatedAt: timestamppb.Now(),
	}, nil
}

func (s *server) GetLong(ctx context.Context, req *desc.GetLongRequest) (*desc.GetLongResponse, error) {

	return &desc.GetLongResponse{
		Long_Url:  "longUrl_" + req.GetShort_Url(), // Или что-то другое
		CreatedAt: timestamppb.Now(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterLongShortV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
