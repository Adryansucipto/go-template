package impl

import (
	"context"
	"go-template/internal/controller/model"
	partnerHTTP "go-template/internal/outbound/http/partner"
	resModel "go-template/pkg/resource/model"
	"go-template/pkg/response"
	"net/http"
)

const serviceName = "partner_service"

type Outbound struct {
	HTTPClient resModel.HTTPClient
}

func New(httpClient resModel.HTTPClient) partnerHTTP.Outbound {
	return &Outbound{
		HTTPClient: httpClient,
	}
}

type partnerUsersResponse struct {
	ResponseCode int          `json:"response_code"`
	Data         []model.User `json:"data"`
}

func (o *Outbound) FetchUsersViaHTTP(ctx context.Context) ([]model.User, error) {
	var resp partnerUsersResponse
	_, err := o.HTTPClient.DoJSON(ctx, resModel.HTTPRequestOption{
		ServiceName: serviceName,
		Method:      http.MethodGet,
		Path:        "/v1/internal/users",
	}, &resp)
	if err != nil {
		return nil, response.ErrorGenerate("[PartnerHTTPOutbound]", err)
	}
	return resp.Data, nil
}
