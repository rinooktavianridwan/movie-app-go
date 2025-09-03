package utils

type StandardResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Version    string      `json:"version"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error,omitempty"`
}

func SuccessResponse(statusCode int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		StatusCode: statusCode,
		Message:    message,
		Version:    "1.0.0",
		Data:       data,
	}
}

func ErrorResponse(statusCode int, error string) StandardResponse {
	return StandardResponse{
		StatusCode: statusCode,
		Message:    error,
		Version:    "1.0.0",
		Data:       nil,
	}
}

func BadRequestResponse(error string) StandardResponse {
	return ErrorResponse(400,  error)
}

func UnauthorizedResponse(error string) StandardResponse {
	return ErrorResponse(401, error)
}

func ForbiddenResponse(error string) StandardResponse {
	return ErrorResponse(403, error)
}

func NotFoundResponse(error string) StandardResponse {
	return ErrorResponse(404, error)
}

func InternalServerErrorResponse(error string) StandardResponse {
	return ErrorResponse(500, error)
}
