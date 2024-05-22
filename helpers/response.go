package helpers

type errorResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

func NewErrorResponse(code int, status string, error string) errorResponse {
	return errorResponse{Code: code, Status: status, Error: error}
}

type successResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessResponse(code int, status string, message string) successResponse {
	return successResponse{Code: code, Status: status, Message: message}
}

type getResponse struct {
	Code     int         `json:"code"`
	Status   string      `json:"status"`
	RowCount int         `json:"rowCount"`
	Data     interface{} `json:"data"`
}

func NewGetResponse(code int, status string, rowCount int, data interface{}) getResponse {
	return getResponse{Code: code, Status: status, RowCount: rowCount, Data: data}
}

type loginResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewLoginResponse(code int, status string, message string) loginResponse {
	return loginResponse{Code: code, Status: status, Message: message}
}

type signUpResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func NewSignUpResponse(code int, status string, message string, token string) signUpResponse {
	return signUpResponse{Code: code, Status: status, Message: message, Token: token}
}