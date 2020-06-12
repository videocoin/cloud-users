package rpc

import (
	"context"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/users/v1"
	ds "github.com/videocoin/cloud-users/datastore"
)

func (s *Server) CreateApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*v1.CreateUserApiTokenResponse, error) { //nolint
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("name", req.Name)

	user, ctx, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	token, err := s.createToken(ctx, user, v1.TokenTypeAPI)
	if err != nil {
		logger.Errorf("failed to create api token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	apiToken, err := s.ds.Token.Create(ctx, user.ID, req.Name, token)
	if err != nil {
		logger.Errorf("failed to create api token record: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &v1.CreateUserApiTokenResponse{
		Id:    apiToken.ID,
		Name:  apiToken.Name,
		Token: apiToken.Token,
	}, nil
}

func (s *Server) ListApiTokens(ctx context.Context, req *protoempty.Empty) (*v1.UserApiListResponse, error) { //nolint
	user, _, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	tokens, err := s.ds.Token.ListByUser(ctx, user.ID)
	if err != nil {
		logger.Errorf("failed to list token by user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	tokensResponse := []*v1.UserApiTokenResponse{}
	if err = copier.Copy(&tokensResponse, tokens); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return &v1.UserApiListResponse{
		Items: tokensResponse,
	}, nil
}

func (s *Server) DeleteApiToken(ctx context.Context, req *v1.UserApiTokenRequest) (*protoempty.Empty, error) { //nolint
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.Id)

	_, ctx, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	err = s.ds.Token.Delete(ctx, req.Id)
	if err != nil {
		s.logger.Errorf("failed to delete api token record: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return new(protoempty.Empty), nil
}

func (s *Server) GetApiToken(ctx context.Context, req *v1.ApiTokenRequest) (*v1.UserApiTokenResponse, error) { //nolint
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("token", req.Token)

	if req.Token == "" {
		return nil, rpc.ErrRpcNotFound
	}

	token, err := s.ds.Token.GetByToken(ctx, req.Token)
	if err != nil {
		if err == ds.ErrTokenNotFound {
			return nil, rpc.ErrRpcNotFound
		}
		return nil, rpc.ErrRpcInternal
	}

	return &v1.UserApiTokenResponse{
		Id:   token.ID,
		Name: token.Name,
	}, nil
}
