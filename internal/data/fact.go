package data

import (
	"sync"

	"aim_testwork/internal/logger"
)

// DTO of Fact
// contains only string values
type FactFormDTO struct {
	PeriodStart string
	PeriodEnd   string
	PeriodKey   string

	IndicatorToMoId     string
	IndicatorToMoFactId string

	Value    string
	FactTime string
	IsPlan   string
	AuthUser string
	Comment  string
}

// Generate all DTO object and send it to chan
func Gen(count int, dtos chan FactFormDTO, pwg *sync.WaitGroup) {
	log := logger.Get()
	log.Info().Msg("generate start")

	for i := 0; i < count; i++ {
		dto := genDto()
		dtos <- dto
	}

	log.Info().Msg("generate completed")
	pwg.Done()
	close(dtos)
}

// Generate single DTO object
func genDto() FactFormDTO {
	return FactFormDTO{
		PeriodStart: "2024-05-01",
		PeriodEnd:   "2024-05-31",
		PeriodKey:   "month",

		IndicatorToMoId:     "227373",
		IndicatorToMoFactId: "0",

		Value:    "1",
		FactTime: "2024-05-31",
		IsPlan:   "0",
		AuthUser: "40",
		Comment:  "Kovalyshev",
	}
}
