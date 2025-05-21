package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"log"
	"log/slog"
	"ppAuthService/internal/config"
	"ppAuthService/internal/repository"
	repoDto "ppAuthService/internal/repository/dto"
	repoErr "ppAuthService/internal/repository/err"
	svcErr "ppAuthService/internal/service/err"
	"ppAuthService/pkg/jwt"
	"ppAuthService/pkg/secure"
	"time"

	proto "github.com/MedvedevEA/ppProtos/gen/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	proto.UnimplementedAuthServiceServer
	store           repository.Repository
	privateKey      *rsa.PrivateKey
	accessLifetime  time.Duration
	refrashLifetime time.Duration
	lg              *slog.Logger
}

func MustNew(store repository.Repository, lg *slog.Logger, cfg *config.Token) *Service {
	privateKey, err := secure.LoadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatalf("failed to initialize service: %v\n", err)
	}

	return &Service{
		store:           store,
		privateKey:      privateKey,
		accessLifetime:  cfg.AccessLifetime,
		refrashLifetime: cfg.RefreshLifetime,
		lg:              lg,
	}
}

func (s *Service) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	userId, err := s.store.AddUser(&repoDto.AddUser{
		Login:    req.Login,
		Password: secure.GetHash(req.Password),
	})
	if err != nil {
		if errors.Is(err, repoErr.ErrUniqueViolation) {
			return nil, status.Error(codes.AlreadyExists, svcErr.ErrLoginAlreadyExists.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.RegisterResponse{UserId: userId.String()}, nil
}
func (s *Service) Unregister(ctx context.Context, req *proto.UnregisterRequest) (*proto.UnregisterResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Unregister"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentUserId.Error())
	}
	if err := s.store.RemoveUser(&userId); err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, svcErr.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.UnregisterResponse{}, nil

}
func (s *Service) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	if req.DeviceCode == "" {
		s.lg.Error("invalid device code value", slog.String("owner", "service.Login"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentDeviceCode.Error())
	}
	user, err := s.store.GetUserByLogin(req.Login)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, svcErr.ErrInvalidLoginOrPassword.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !secure.CheckHash(req.Password, user.Password) {
		s.lg.Error("hash verification error", slog.String("owner", "service.Login"))
		return nil, status.Error(codes.Unauthenticated, svcErr.ErrInvalidLoginOrPassword.Error())
	}
	if err := s.store.RevokeRefreshTokensByUserIdAndDeviceCode(&repoDto.RevokeRefreshTokensByUserIdAndDeviceCode{
		UserId:     user.UserId,
		DeviceCode: &req.DeviceCode,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	//access token
	accessTokenString, _, err := jwt.CreateToken(user.UserId, req.DeviceCode, "access", s.accessLifetime, s.privateKey)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Login"))
		return nil, status.Error(codes.Internal, err.Error())
	}
	//refresh token
	refreshTokenString, refreshTokenClaims, err := jwt.CreateToken(user.UserId, req.DeviceCode, "refresh", s.refrashLifetime, s.privateKey)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Login"))
	}
	if err := s.store.AddRefreshTokenWithRefreshTokenId(&repoDto.AddRefreshTokenWithRefreshTokenId{
		RefreshTokenId: refreshTokenClaims.Jti,
		UserId:         refreshTokenClaims.Sub,
		DeviceCode:     refreshTokenClaims.DeviceCode,
		ExpirationAt:   refreshTokenClaims.ExpiresAt.Time,
		IsRevoke:       false,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.LoginResponse{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}
func (s *Service) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.Logout"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentUserId.Error())
	}
	if req.DeviceCode == "" {
		s.lg.Error("invalid device code value", slog.String("owner", "service.Logout"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentDeviceCode.Error())
	}
	if err := s.store.RevokeRefreshTokensByUserIdAndDeviceCode(&repoDto.RevokeRefreshTokensByUserIdAndDeviceCode{
		UserId:     &userId,
		DeviceCode: &req.DeviceCode,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.LogoutResponse{}, nil
}
func (s *Service) UpdatePassword(ctx context.Context, req *proto.UpdatePasswordRequest) (*proto.UpdatePasswordResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.UpdatePassword"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentUserId.Error())
	}
	hashNewPassword := secure.GetHash(req.NewPassword)
	if err = s.store.UpdateUser(&repoDto.UpdateUser{
		UserId:   &userId,
		Password: &hashNewPassword,
	}); err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, svcErr.ErrUserNotFound.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := s.store.RevokeRefreshTokensByUserIdAndDeviceCode(&repoDto.RevokeRefreshTokensByUserIdAndDeviceCode{
		UserId:     &userId,
		DeviceCode: nil,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.UpdatePasswordResponse{}, nil
}
func (s *Service) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	refreshTokenId, err := uuid.Parse(req.RefreshTokenId)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RefreshToken"))
		return nil, status.Error(codes.InvalidArgument, svcErr.ErrInvalidArgumentTokenId.Error())
	}
	refreshToken, err := s.store.GetRefreshToken(&refreshTokenId)
	if err != nil {
		if errors.Is(err, repoErr.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, svcErr.ErrTokenNotFound.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if refreshToken.IsRevoke {
		if err := s.store.RevokeRefreshTokensByUserIdAndDeviceCode(&repoDto.RevokeRefreshTokensByUserIdAndDeviceCode{
			UserId:     refreshToken.UserId,
			DeviceCode: &refreshToken.DeviceCode,
		}); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		s.lg.Info("token is revoked", slog.String("owner", "service.RefreshToken"), slog.Any("tokenId", refreshToken.RefreshTokenId))
		return nil, status.Error(codes.Unauthenticated, svcErr.ErrTokenRevoked.Error())
	}
	if err := s.store.RevokeRefreshTokenByRefreshTokenId(&refreshTokenId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	//access token
	accessTokenString, _, err := jwt.CreateToken(refreshToken.UserId, refreshToken.DeviceCode, "access", s.accessLifetime, s.privateKey)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RefreshToken"))
		return nil, status.Error(codes.Internal, err.Error())
	}
	//refresh token
	refreshTokenString, refreshTokenClaims, err := jwt.CreateToken(refreshToken.UserId, refreshToken.DeviceCode, "refresh", s.refrashLifetime, s.privateKey)
	if err != nil {
		s.lg.Error(err.Error(), slog.String("owner", "service.RefreshToken"))
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := s.store.AddRefreshTokenWithRefreshTokenId(&repoDto.AddRefreshTokenWithRefreshTokenId{
		RefreshTokenId: refreshTokenClaims.Jti,
		UserId:         refreshTokenClaims.Sub,
		DeviceCode:     refreshTokenClaims.DeviceCode,
		ExpirationAt:   refreshTokenClaims.ExpiresAt.Time,
		IsRevoke:       false,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.RefreshTokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
