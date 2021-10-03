package dto

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
