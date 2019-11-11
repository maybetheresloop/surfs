package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"surfs/internal/block"

	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Welcome to surfs-block!")
	lis, err := net.Listen("tcp", ":7878")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	store := block.NewStore()
	s := grpc.NewServer()
	block.RegisterStoreServer(s, store)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
