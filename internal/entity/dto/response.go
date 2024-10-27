package dto

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseDetail struct {
	Detail string `json:"detail"`
}
