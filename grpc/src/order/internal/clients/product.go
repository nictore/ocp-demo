package clients

import (
	"log"
	"time"

	"github.com/iamdpastore/ocp-demo/grpc/src/proto/product"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var ProductService product.ProductServiceClient

func init() {
	log.Printf("[Order] Dialing to 'product:5000' ...")

	keepAliveParams := keepalive.ClientParameters{
		Time:                5 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.Dial("product:5000", grpc.WithInsecure(), grpc.WithBlock(), grpc.FailOnNonTempDialError(true), grpc.WithKeepaliveParams(keepAliveParams), grpc.WithStreamInterceptor(
		grpc_opentracing.StreamClientInterceptor(
			grpc_opentracing.WithTracer(opentracing.GlobalTracer()))))
	if err != nil {
		log.Fatalf("[Order] Error dialing to Product service: %v", err)
	}

	ProductService = product.NewProductServiceClient(conn)
}
