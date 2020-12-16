package model

type ServiceWrapperVo struct {
	Data []ServiceVo `json:"data"`
}

type CheckWrapperVo struct {
	Data []CheckVo `json:"data"`
}

type FailureWrapperVo struct {
	Data []FailureVo `json:"data"`
}
