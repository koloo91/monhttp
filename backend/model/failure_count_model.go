package model

type FailureCount struct {
	Count int
}

type FailureCountVo struct {
	Count int `json:"count"`
}

func MapFailureCountEntityToVo(entity FailureCount) FailureCountVo {
	return FailureCountVo{Count: entity.Count}
}
