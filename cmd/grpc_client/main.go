package main

import (
	"context"
	"log"
	"time"

	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	desc "github.com/MariaPopova99/microservices/pkg/note_v1"
)

const (
	address = "localhost:50051"
)

func getLongUrl(c desc.LongShortV1Client, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetLong(ctx, &desc.GetLongRequest{ShortUrl: url})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}
	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetLongUrl()))
}

func getShortUrl(c desc.LongShortV1Client, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetShort(ctx, &desc.GetShortRequest{LongUrl: url})
	if err != nil {
		log.Fatalf("failed to get note by id: %v", err)
	}

	log.Printf(color.RedString("Note info:\n"), color.GreenString("%+v", r.GetShortUrl()))
}

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	c := desc.NewLongShortV1Client(conn)

	getLongUrl(c, "longUrlToConvert")
	getShortUrl(c, "ShortUrlToConvert")
	getLongUrl(c, "e9ddaa35")
	getShortUrl(c, "MyNewSiteLongURL4") //new generation
	getLongUrl(c, "ssadsgsge")          //no rows error

}
