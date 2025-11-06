package dto

type HttpResponse[TData any] struct {
	HttpStatus int    `json:"http_status"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Data       TData  `json:"repository"`
}

type Page[T any] struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
	Items []T `json:"items"`
}
