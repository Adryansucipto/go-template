package consts

const (
	INTERNAL_SERVER_ERROR_CODE = 500
	BAD_REQUEST_CODE           = 400
	NOT_FOUND_CODE             = 404
	SUCCESS_CODE               = 200
	UNAUTHORIZED               = 401
)

var ErrorMessages = map[int]string{
	INTERNAL_SERVER_ERROR_CODE: "Internal Server Error",
	BAD_REQUEST_CODE:           "Bad Request",
	NOT_FOUND_CODE:             "Not Found",
	SUCCESS_CODE:               "Success",
	UNAUTHORIZED:               "Unauthorized",
}
