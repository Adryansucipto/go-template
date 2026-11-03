package impl

import (
	"context"
	"go-template/internal/controller/model"
	partnerHTTPRepo "go-template/internal/repository/client/http/partner"
	resModel "go-template/pkg/resource/model"
	"go-template/pkg/response"
	"net/http"
)

const serviceName = "partner_service"

type Repository struct {
	HTTPClient resModel.HTTPClient
}

func New(httpClient resModel.HTTPClient) partnerHTTPRepo.Repository {
	return &Repository{
		HTTPClient: httpClient,
	}
}

type partnerUsersResponse struct {
	ResponseCode int          `json:"response_code"`
	Data         []model.User `json:"data"`
}

func (r *Repository) FetchUsersViaHTTP(ctx context.Context) ([]model.User, error) {
	var resp partnerUsersResponse
	_, err := r.HTTPClient.DoJSON(ctx, resModel.HTTPRequestOption{
		ServiceName: serviceName,
		Method:      http.MethodGet,
		Path:        "/v1/internal/users",
	}, &resp)
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerHTTPClient]", err)
	}
	return resp.Data, nil
}

