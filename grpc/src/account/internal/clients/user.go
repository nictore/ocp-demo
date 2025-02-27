package clients

import (
	"log"
	"time"

	"github.com/iamdpastore/ocp-demo/grpc/src/proto/user"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var UserService user.UserServiceClient

func init() {
	log.Printf("[Account] Dialing to 'user:5000' ...")

	keepAliveParams := keepalive.ClientParameters{
		Time:                5 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.Dial("user:5000", grpc.WithInsecure(), grpc.WithBlock(), grpc.FailOnNonTempDialError(true), grpc.WithKeepaliveParams(keepAliveParams), grpc.WithStreamInterceptor(
		grpc_opentracing.StreamClientInterceptor(
			grpc_opentracing.WithTracer(opentracing.GlobalTracer()))))
	if err != nil {
		log.Fatalf("[Account] Error dialing to User service: %v", err)
	}

	UserService = user.NewUserServiceClient(conn)
}
