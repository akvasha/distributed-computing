package main

import (
	"context"
	"google.golang.org/grpc"
	pbauth "lib/proto"
	"log"
	"net"
)

type AuthRpcServer struct {
}

func (s *AuthRpcServer) Validate(c context.Context, req *pbauth.ValidateRequest) (resp *pbauth.ValidateResponse, err error) {
	var permission UserPermission
	if permission, err = ValidateAccess(req.AccessToken); err != nil {
		return
	}
	resp = &pbauth.ValidateResponse{
		Username: permission.Username,
		Admin:    permission.Admin,
	}
	return
}

func RunRPCServer() {
	listener, err := net.Listen("tcp", ":5300")
	if err != nil {
		log.Fatal("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	pbauth.RegisterAuthRpcServer(grpcServer, &AuthRpcServer{})
	log.Println("RPC server started")
	log.Fatal(grpcServer.Serve(listener))
}
