package util

import "github.com/labstack/echo/v4"

type Response struct {
	ResponseCode  int   `json:"response_code"`
	ResponseError error `json:"response_error"`
	Data          any   `json:"data,omitempty"`
	Page          any   `json:"page,omitempty"`
}

func ResponseGenerate(code int, err error, data any, page any) Response {
	return Response{
		ResponseCode:  code,
		ResponseError: err,
		Data:          data,
		Page:          page,
	}
}

func HttpResponses(eCtx echo.Context, response Response) error {
	return eCtx.JSON(response.ResponseCode, response)
}
