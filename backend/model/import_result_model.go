package model

type ImportResult struct {
	RowNumber int
	Service   Service
	Error     error
}

type ImportResultVo struct {
	RowNumber int       `json:"rowNumber"`
	Service   ServiceVo `json:"service"`
	Error     string    `json:"error"`
}

func MapImportResultEntityToVo(entity ImportResult) ImportResultVo {
	errorString := ""
	if entity.Error != nil {
		errorString = entity.Error.Error()
	}

	return ImportResultVo{
		RowNumber: entity.RowNumber,
		Service:   MapServiceEntityToVo(entity.Service),
		Error:     errorString,
	}
}

func MapImportResultEntitiesToVos(entities []ImportResult) []ImportResultVo {
	result := make([]ImportResultVo, 0, len(entities))

	for _, entity := range entities {
		result = append(result, MapImportResultEntityToVo(entity))
	}

	return result
}
