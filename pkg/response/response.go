package response

import (
	"go-template/pkg/response/consts"

	"github.com/labstack/echo/v4"
)

type Response struct {
	ResponseCode    int   `json:"response_code"`
	ResponseMessage error `json:"response_message"`
	Errors          any   `json:"errors,omitempty"`
	Data            any   `json:"data,omitempty"`
	Page            any   `json:"page,omitempty"`
}

func (r Response) Error() string {
	if r.ResponseMessage != nil {
		return r.ResponseMessage.Error()
	}
	return consts.ErrorMessages[r.ResponseCode]
}

type FinalResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Errors          any    `json:"errors,omitempty"`
	Data            any    `json:"data,omitempty"`
	Page            any    `json:"page,omitempty"`
}

func ResponseGenerate(code int, msg error, errors any, data any, page any) Response {
	return Response{
		ResponseCode:    code,
		ResponseMessage: msg,
		Errors:          errors,
		Data:            data,
		Page:            page,
	}
}

func HttpResponses(eCtx echo.Context, res Response) error {
	var finalResponse FinalResponse
	finalResponse.ResponseCode = res.ResponseCode
	finalResponse.ResponseMessage = consts.ErrorMessages[res.ResponseCode]
	finalResponse.Data = res.Data
	finalResponse.Errors = res.Errors
	finalResponse.Page = res.Page
	return eCtx.JSON(res.ResponseCode, finalResponse)
}
