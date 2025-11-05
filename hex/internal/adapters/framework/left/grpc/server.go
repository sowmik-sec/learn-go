package grpc

import (
	"hex/internal/adapters/framework/left/grpc/pb"
	"hex/internal/ports"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api ports.APIPort
	pb.UnimplementedArithmeticServiceServer
}

func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (grpca Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", "127.0.0.1:9001")
	if err != nil {
		log.Fatalf("failed to listen on 127.0.0.1:9001: %v", err)
	}
	log.Println("Listening on 127.0.0.1:9001")
	arithmeticServiceServer := grpca
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterArithmeticServiceServer(grpcServer, arithmeticServiceServer)
	log.Println("Starting gRPC server on port 9001")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port 9000: %v", err)
	}
}
