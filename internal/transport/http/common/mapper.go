package common

func ToErrorResponse(err string) ErrorResponse {
	return ErrorResponse{Error: err}
}

func ToSuccessfulResponse(msg string) SuccessfulResponse {
	return SuccessfulResponse{Message: msg}
}
