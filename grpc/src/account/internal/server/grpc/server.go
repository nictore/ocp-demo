package grpc

import (
	"log"
	"net"
	"sync"

	"github.com/iamdpastore/ocp-demo/grpc/src/account/internal/impl"
	"github.com/iamdpastore/ocp-demo/grpc/src/proto/account"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func Serve(wg *sync.WaitGroup, port string) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("[Account] GRPC failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
	)))

	account.RegisterAccountServiceServer(s, &impl.Server{})

	log.Printf("[Account] Serving GRPC on localhost:%s ...", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[Account] GRPC failed to serve: %v", err)
	}
}
