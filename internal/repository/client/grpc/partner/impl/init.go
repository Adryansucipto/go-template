package impl

import (
	"context"
	partnerGRPCRepo "go-template/internal/repository/client/grpc/partner"
	resModel "go-template/pkg/resource/model"
	"go-template/pkg/response"
	userv1 "go-template/proto/user/v1"
)

const serviceName = "partner_service"

type Repository struct {
	GRPCClient *resModel.GRPCClient
}

func New(grpcClient *resModel.GRPCClient) partnerGRPCRepo.Repository {
	return &Repository{
		GRPCClient: grpcClient,
	}
}

func (r *Repository) FetchUsersViaGRPC(ctx context.Context) ([]*userv1.UserItem, error) {
	conn, err := r.GRPCClient.GetConnection(serviceName)
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerGRPCClient]", err)
	}

	client := userv1.NewUserServiceClient(conn)
	resp, err := client.GetUsers(ctx, &userv1.GetUsersRequest{})
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerGRPCClient]", err)
	}

	return resp.Users, nil
}

