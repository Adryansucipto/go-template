package impl

import (
	"context"
	partnerGRPC "go-template/internal/outbound/grpc/partner"
	resModel "go-template/pkg/resource/model"
	"go-template/pkg/response"
	userv1 "go-template/proto/user/v1"
)

const serviceName = "partner_service"

type Outbound struct {
	GRPCClient *resModel.GRPCClient
}

func New(grpcClient *resModel.GRPCClient) partnerGRPC.Outbound {
	return &Outbound{
		GRPCClient: grpcClient,
	}
}

func (o *Outbound) FetchUsersViaGRPC(ctx context.Context) ([]*userv1.UserItem, error) {
	conn, err := o.GRPCClient.GetConnection(serviceName)
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerGRPCOutbound]", err)
	}

	client := userv1.NewUserServiceClient(conn)
	resp, err := client.GetUsers(ctx, &userv1.GetUsersRequest{})
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerGRPCOutbound]", err)
	}

	return resp.Users, nil
}
