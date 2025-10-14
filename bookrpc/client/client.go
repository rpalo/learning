package main

import (
	"context"
	"log"
	"time"

	"github.com/rpalo/learning/bookrpc/books"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Couldn't connect to server: %v", err)
	}
	defer conn.Close()

	client := books.NewBooksClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.CreateList(ctx, &books.CreateListRequest{
		Name: "jeremiah",
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("fofof: %v", resp)
	req := &books.GetListRequest{Id: resp.Id}
	list, err := client.GetList(ctx, req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(list)
}
