package model

type Average struct {
	LastDayResponseTime  int
	LastWeekResponseTime int
	LastDayUptime        float64
	LastWeekUptime       float64
}

type AverageVo struct {
	LastDayResponseTime  int     `json:"lastDayResponseTime"`
	LastWeekResponseTime int     `json:"lastWeekResponseTime"`
	LastDayUptime        float64 `json:"lastDayUptime"`
	LastWeekUptime       float64 `json:"lastWeekUptime"`
}

func MapAverageEntityToVo(entity Average) AverageVo {
	return AverageVo{
		LastDayResponseTime:  entity.LastDayResponseTime,
		LastWeekResponseTime: entity.LastWeekResponseTime,
		LastDayUptime:        entity.LastDayUptime,
		LastWeekUptime:       entity.LastWeekUptime,
	}
}
