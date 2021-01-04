package model

import "time"

type FailureCountByDay struct {
	Day   time.Time
	Count int
}

type FailureCountByDayVo struct {
	Day   time.Time `json:"day"`
	Count int       `json:"count"`
}

func MapFailureCountByDayEntityToVo(entity FailureCountByDay) FailureCountByDayVo {
	return FailureCountByDayVo{
		Day:   entity.Day,
		Count: entity.Count,
	}
}

func MapFailureCountByDayEntitiesToVos(entities []FailureCountByDay) []FailureCountByDayVo {
	result := make([]FailureCountByDayVo, 0, len(entities))

	for _, entity := range entities {
		result = append(result, MapFailureCountByDayEntityToVo(entity))
	}
	return result
}
