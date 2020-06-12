package rpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRPCUserAlreadyExists = status.Errorf(codes.AlreadyExists, "User already exists")
)
