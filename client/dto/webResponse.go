package dto

type WebSuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type WebErrorResponse struct {
	Status int         `json:"status"`
	Type   string      `json:"type"`
	Detail interface{} `json:"detail"`
}
