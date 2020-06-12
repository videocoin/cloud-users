package rpc

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	protoempty "github.com/gogo/protobuf/types"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	accountsv1 "github.com/videocoin/cloud-api/accounts/v1"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/users/v1"
	"github.com/videocoin/cloud-pkg/auth"
	ds "github.com/videocoin/cloud-users/datastore"
)

func (s *Server) Validate(ctx context.Context, req *v1.ValidateUserRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	logger := s.logger.WithField("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	err := s.ds.User.Validate(ctx, req.Email)
	if err != nil {
		if err == ds.ErrUserAlreadyExists {
			respErr := &rpc.MultiValidationError{
				Errors: []*rpc.ValidationError{
					{
						Field:   "email",
						Message: "Email is already registered",
					},
				},
			}
			return nil, rpc.NewRpcValidationError(respErr)
		}

		logger.Error(err)

		return nil, rpc.ErrRpcInternal
	}

	return new(protoempty.Empty), nil
}

func (s *Server) Create(ctx context.Context, req *v1.CreateUserRequest) (*v1.TokenResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	logger := s.logger.WithField("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.Register(ctx, req)
	if err != nil {
		if err == ds.ErrUserAlreadyExists {
			respErr := &rpc.MultiValidationError{
				Errors: []*rpc.ValidationError{
					{
						Field:   "email",
						Message: "Email is already registered",
					},
				},
			}
			return nil, rpc.NewRpcValidationError(respErr)
		}

		logger.Error(err)

		return nil, rpc.ErrRpcInternal
	}

	logger = logger.WithField("id", user.ID)

	token, err := s.createToken(ctx, user, v1.TokenTypeRegular)
	if err != nil {
		logger.Errorf("failed to create token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.UpdateAuthToken(ctx, user, token); err != nil {
		logger.Errorf("failed to update auth token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	accReq := &accountsv1.AccountRequest{OwnerId: user.ID}
	if err = s.eb.CreateUserAccount(span, accReq); err != nil {
		logger.Errorf("failed to create account via eventbus: %s", err)
	}

	confirmToken := newRecoveryToken(req.Email, 1*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))
	if err = s.notifications.SendEmailConfirmation(ctx, user, confirmToken); err != nil {
		logger.WithField("failed to send confirm email to user id", user.ID).Error(err)
	}

	return &v1.TokenResponse{
		Token: token,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UserProfile, error) {
	user, _, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	if err = s.ds.User.UpdateUIRole(ctx, user, req.UiRole); err != nil {
		logger.Errorf("failed to update user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return userProfile, nil
}

func (s *Server) Login(ctx context.Context, req *v1.LoginUserRequest) (*v1.TokenResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	logger := s.logger.WithField("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == ds.ErrUserNotFound {
			return nil, rpc.ErrRpcUnauthenticated
		}

		logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	logger = logger.WithField("id", user.ID)

	if !checkPasswordHash(ctx, req.Password, user.Password) {
		return nil, rpc.ErrRpcUnauthenticated
	}

	token, err := s.createToken(ctx, user, v1.TokenTypeRegular)
	if err != nil {
		logger.Errorf("failed to create token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.UpdateAuthToken(ctx, user, token); err != nil {
		logger.Errorf("failed to update auth token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &v1.TokenResponse{
		Token: user.Token,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
	user, _, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	if err = s.ds.User.ResetAuthToken(user); err != nil {
		logger.Errorf("failed to reset auth token: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *Server) StartRecovery(ctx context.Context, req *v1.StartRecoveryUserRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("email", req.Email)

	logger := s.logger.WithField("email", req.Email)

	if verr := s.validator.validate(req); verr != nil {
		logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := s.ds.User.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == ds.ErrUserNotFound {
			return nil, rpc.ErrRpcBadRequest
		}

		logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	token := newRecoveryToken(req.Email, 12*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))
	if err = s.notifications.SendEmailRecovery(ctx, user, token); err != nil {
		logger.WithField("failed to send recovery email to user id", user.ID).Error(err)
	}

	return &protoempty.Empty{}, nil
}

func (s *Server) Recover(ctx context.Context, req *v1.RecoverUserRequest) (*protoempty.Empty, error) {
	if err := s.validator.validate(req); err != nil {
		s.logger.Warning(err)
		return nil, rpc.NewRpcValidationError(err)
	}

	user, err := verifyRecoveryToken(ctx, req.Token, s.ds.User.GetByEmail, []byte(s.authRecoverySecret))
	if err != nil {
		if err == ds.ErrUserNotFound {
			return nil, rpc.ErrRpcBadRequest
		}

		s.logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	if err = s.ds.User.ResetPassword(ctx, user, req.Password); err != nil {
		s.logger.Errorf("failed to reset password: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil

}

func (s *Server) StartConfirmation(ctx context.Context, req *protoempty.Empty) (*protoempty.Empty, error) {
	user, _, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	confirmToken := newRecoveryToken(user.Email, 1*time.Hour, []byte(user.Password), []byte(s.authRecoverySecret))
	if err = s.notifications.SendEmailConfirmation(ctx, user, confirmToken); err != nil {
		logger.
			WithField("user_id", user.ID).
			WithError(err).
			Error("failed to send confirmation email to user")
	}

	return &protoempty.Empty{}, nil
}

func (s *Server) Confirm(ctx context.Context, req *v1.ConfirmUserRequest) (*protoempty.Empty, error) {
	if verr := s.validator.validate(req); verr != nil {
		s.logger.Warning(verr)
		return nil, rpc.NewRpcValidationError(verr)
	}

	user, err := verifyRecoveryToken(ctx, req.Token, s.ds.User.GetByEmail, []byte(s.authRecoverySecret))
	if err != nil {
		if err == ds.ErrUserNotFound {
			return nil, rpc.ErrRpcBadRequest
		}

		s.logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	logger := s.logger.WithField("id", user.ID)

	if err := s.ds.User.Activate(user.ID); err != nil {
		logger.Errorf("failed to activate user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *v1.ResetPasswordUserRequest) (*protoempty.Empty, error) {
	user, _, err := s.authenticate(ctx, false)
	if err != nil {
		return nil, err
	}

	logger := s.logger.WithField("id", user.ID)

	if err = s.ds.User.ResetPassword(ctx, user, req.Password); err != nil {
		logger.Errorf("failed to reset password: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	return &protoempty.Empty{}, nil
}

func (s *Server) Get(ctx context.Context, req *protoempty.Empty) (*v1.UserProfile, error) {
	user, _, err := s.authenticate(ctx, true)
	if err != nil {
		return nil, err
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return userProfile, nil
}

func (s *Server) GetById(ctx context.Context, req *v1.UserRequest) (*v1.UserProfile, error) { //nolint
	user, err := s.ds.User.Get(req.Id)
	if err != nil {
		if err == ds.ErrUserNotFound {
			return nil, rpc.ErrRpcNotFound
		}

		s.logger.Errorf("failed to get user: %s", err)
		return nil, rpc.ErrRpcInternal
	}

	userProfile := new(v1.UserProfile)
	if err = copier.Copy(userProfile, user); err != nil {
		return nil, rpc.ErrRpcInternal
	}

	return userProfile, nil
}

func (s *Server) createToken(ctx context.Context, user *ds.User, tokenType v1.TokenType) (string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "createToken")
	defer span.Finish()

	span.SetTag("id", user.ID)
	span.SetTag("email", user.Email)
	span.SetTag("token_type", tokenType)

	claims := auth.ExtendedClaims{
		Type: auth.TokenType(tokenType),
		StandardClaims: jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	if tokenType == v1.TokenTypeAPI {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Hour * 24 * 365).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.authTokenSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s *Server) authenticate(ctx context.Context, allowAPI bool) (*ds.User, context.Context, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "authenticate")
	defer span.Finish()

	ctx = auth.NewContextWithSecretKey(ctx, s.authTokenSecret)
	ctx, _, err := auth.AuthFromContext(ctx)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	if !allowAPI && s.getTokenType(ctx) == auth.TokenType(v1.TokenTypeAPI) {
		return nil, nil, rpc.ErrRpcPermissionDenied
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	user, err := s.ds.User.Get(userID)
	if err != nil {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	if user.Token == "" {
		return nil, ctx, rpc.ErrRpcUnauthenticated
	}

	return user, ctx, nil
}

func (s *Server) getTokenType(ctx context.Context) auth.TokenType {
	tokenType, ok := auth.TypeFromContext(ctx)
	if !ok {
		return auth.TokenType(v1.TokenTypeRegular)
	}

	return tokenType
}
