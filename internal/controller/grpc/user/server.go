package user

import (
	"context"
	"go-template/internal/controller/model"
	userUsecase "go-template/internal/usecase/user"
	userv1 "go-template/proto/user/v1"

	"go.uber.org/dig"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	dig.In

	User   userUsecase.Usecase
	Logger *zap.Logger
}

func (s *Server) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	res, err := s.User.GetUsers(ctx)
	if err != nil {
		s.Logger.Error("[UserGRPCServer] GetUsers failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	users, ok := res.Data.([]model.User)
	if !ok {
		return &userv1.GetUsersResponse{}, nil
	}

	var items []*userv1.UserItem
	for _, u := range users {
		items = append(items, &userv1.UserItem{
			ID:       int32(u.ID),
			Name:     u.Name,
			Username: u.Username,
		})
	}

	return &userv1.GetUsersResponse{
		Users: items,
	}, nil
}

