package utils

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewErrorResponse(error, message string) errorResponse {
	return errorResponse{Error: error, Message: message}
}

type successResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessResponse(status, message string) successResponse {
	return successResponse{Status: status, Message: message}
}

type getResponse struct {
	RowCount int         `json:"rowCount"`
	Data     interface{} `json:"data"`
}

func NewGetResponse(rowCount int, data interface{}) getResponse {
	return getResponse{RowCount: rowCount, Data: data}
}

type loginResponse struct {
	Message string `json:"message"`
}

func NewLoginResponse(message string) loginResponse {
	return loginResponse{Message: message}
}

type signUpResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func NewSignUpResponse(message string, token string) signUpResponse {
	return signUpResponse{Message: message, Token: token}
}
