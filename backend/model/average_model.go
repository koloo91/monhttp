package model

type Average struct {
	LastDay  int
	LastWeek int
}

type AverageVo struct {
	LastDay  int `json:"lastDay"`
	LastWeek int `json:"lastWeek"`
}

func MapAverageEntityToVo(entity Average) AverageVo {
	return AverageVo{
		LastDay:  entity.LastDay,
		LastWeek: entity.LastWeek,
	}
}
