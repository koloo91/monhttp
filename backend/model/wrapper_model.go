package model

type ServiceWrapperVo struct {
	Data       []ServiceVo `json:"data"`
	TotalCount int         `json:"totalCount"`
	PageSize   int         `json:"pageSize"`
	Page       int         `json:"page"`
}

type CheckWrapperVo struct {
	Data []CheckVo `json:"data"`
}

type FailureWrapperVo struct {
	Data       []FailureVo `json:"data"`
	TotalCount int         `json:"totalCount"`
	PageSize   int         `json:"pageSize"`
	Page       int         `json:"page"`
}

type NotifierWrapperVo struct {
	Data []NotifierVo `json:"data"`
}
