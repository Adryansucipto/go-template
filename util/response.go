package util

import (
	"go-template/util/consts"

	"github.com/labstack/echo/v4"
)

type Response struct {
	ResponseCode    int   `json:"response_code"`
	ResponseMessage error `json:"response_message"`
	Errors          any   `json:"errors,omitempty"`
	Data            any   `json:"data,omitempty"`
	Page            any   `json:"page,omitempty"`
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

func HttpResponses(eCtx echo.Context, response Response) error {
	var finalResponse FinalResponse
	finalResponse.ResponseCode = response.ResponseCode
	finalResponse.ResponseMessage = consts.ErrorMessages[response.ResponseCode]
	finalResponse.Data = response.Data
	finalResponse.Errors = response.Errors
	finalResponse.Page = response.Page
	return eCtx.JSON(response.ResponseCode, finalResponse)
}
