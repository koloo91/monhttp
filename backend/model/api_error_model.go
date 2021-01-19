package model

type ApiErrorVo struct {
	Message string        `json:"message"`
	Errors  []interface{} `json:"errors"`
}
