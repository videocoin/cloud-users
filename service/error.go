package service

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var (
	ErrRpcUserAlreadyExists = grpc.Errorf(codes.AlreadyExists, "User already exists")
)
