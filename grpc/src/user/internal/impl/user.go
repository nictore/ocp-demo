package impl

import (
	"context"
	"log"
	"strconv"

	"github.com/Pallinder/go-randomdata"
	"github.com/iamdpastore/ocp-demo/grpc/src/proto/user"
)

type Server struct {
	user.UnimplementedUserServiceServer
}

func (s *Server) Create(ctx context.Context, in *user.CreateUserReq) (*user.CreateUserResp, error) {

	log.Printf("[User] Create Req: %v", in.GetUser())

	r := &user.CreateUserResp{Id: strconv.Itoa(randomdata.Number(1000000))}

	log.Printf("[User] Create Res: %v", r.GetId())

	return r, nil
}

func (s *Server) Read(ctx context.Context, in *user.ReadUserReq) (*user.ReadUserResp, error) {

	log.Printf("[User] Read Req: %v", in.GetId())

	r := &user.ReadUserResp{User: &user.User{Id: in.GetId(), Name: randomdata.FullName(randomdata.RandomGender), Email: randomdata.Email()}}

	log.Printf("[User] Read Res: %v", r.GetUser())

	return r, nil
}
