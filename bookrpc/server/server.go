package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/rpalo/learning/bookrpc/books"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BooksServer struct {
	books.UnimplementedBooksServer
	db map[string]*books.List
}

func NewBooksServer() *BooksServer {
	return &BooksServer{
		db: make(map[string]*books.List),
	}
}

func (server *BooksServer) CreateList(
	ctx context.Context,
	req *books.CreateListRequest,
) (*books.List, error) {
	// if _, ok := server.db[req.Name]; ok {
	// 	return nil, status.Error(codes.AlreadyExists, "list with that name already exists")
	// }
	bookId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	list := &books.List{
		Id:   bookId.String(),
		Name: req.Name,
	}
	server.db[list.Id] = list
	return list, nil
}

func (server *BooksServer) GetList(
	ctx context.Context,
	req *books.GetListRequest,
) (*books.List, error) {
	if list, ok := server.db[req.Id]; ok {
		return list, nil
	} else {
		return nil, status.Error(codes.NotFound, "not found")
	}
}

func main() {
	port := flag.Int("port", 8000, "the port to serve on")
	flag.Parse()
	log.Printf("Starting server on port %d...", *port)
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	books.RegisterBooksServer(grpcServer, NewBooksServer())

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Couldn't set up listener: %v", err)
	}
	grpcServer.Serve(listener)
}
