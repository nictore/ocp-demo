package grpc

import (
	"log"
	"net"
	"sync"

	"github.com/iamdpastore/ocp-demo/grpc/src/proto/user"
	"github.com/iamdpastore/ocp-demo/grpc/src/user/internal/impl"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func Serve(wg *sync.WaitGroup, port string) {
	defer wg.Done()

	lis, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("[User] GRPC failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
	)))

	user.RegisterUserServiceServer(s, &impl.Server{})

	log.Printf("[User] Serving GRPC on localhost:%s ...", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("[User] GRPC failed to serve: %v", err)
	}
}
