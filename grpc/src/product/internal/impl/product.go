package impl

import (
	"context"
	"log"
	"strconv"

	"github.com/Pallinder/go-randomdata"
	"github.com/iamdpastore/ocp-demo/grpc/src/proto/product"
)

type Server struct {
	product.UnimplementedProductServiceServer
}

func (s *Server) Create(ctx context.Context, in *product.CreateProductReq) (*product.CreateProductResp, error) {

	log.Printf("[Product] Create Req: %v", in.GetProduct())

	r := &product.CreateProductResp{Id: strconv.Itoa(randomdata.Number(1000000))}

	log.Printf("[Product] Create Res: %v", r.GetId())

	return r, nil
}

func (s *Server) Read(ctx context.Context, in *product.ReadProductReq) (*product.ReadProductResp, error) {

	log.Printf("[Product] Read Req: %v", in.GetId())

	r := &product.ReadProductResp{Product: &product.Product{Id: in.GetId(), Name: randomdata.SillyName(), Description: randomdata.Paragraph(), Price: int32(randomdata.Number(1000))}}

	log.Printf("[Product] Read Res: %v", r.GetProduct())

	return r, nil
}
