package sender

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"

	"aim_testwork/internal/conf"
	"aim_testwork/internal/data"
	"aim_testwork/internal/logger"
)

// function to start sending information
func SendToServer(dtos chan data.FactFormDTO, pwg *sync.WaitGroup) {
	log := logger.Get()
	log.Info().Msg("sender started")
	i := 0

	for dto := range dtos {
		i++
		log.Info().Int("Index", i).Msg("try to send")
		sendDtoToServer(dto)
	}

	log.Info().Msg("sender comleted")
	pwg.Done()
}

// send single DTO object
func sendDtoToServer(dto data.FactFormDTO) {
	cfg := conf.Get()
	log := logger.Get()
	client := &http.Client{}

	// Format form/data is worked through nginx and gets 200 status code
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("period_start", dto.PeriodStart)
	writer.WriteField("period_end", dto.PeriodEnd)
	writer.WriteField("period_key", dto.PeriodKey)
	writer.WriteField("indicator_to_mo_id", dto.IndicatorToMoId)
	writer.WriteField("indicator_to_mo_fact_id", dto.IndicatorToMoFactId)
	writer.WriteField("value", dto.Value)
	writer.WriteField("fact_time", dto.FactTime)
	writer.WriteField("is_plan", dto.IsPlan)
	writer.WriteField("auth_user_id", dto.AuthUser)
	writer.WriteField("comment", dto.Comment)
	writer.Close()
	log.Debug().Str("Action", "Request").Str("Body", body.String()).Send()

	req, err := http.NewRequest("POST", cfg.Server.Url, bytes.NewReader(body.Bytes()))
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// Format "form-urlencoded" is not worked and gets 500 status code

	// form := url.Values{}
	// form.Set("period_start", dto.PeriodStart)
	// form.Set("period_end", dto.PeriodEnd)
	// form.Set("period_key", dto.PeriodKey)
	// form.Set("indicator_to_mo_id", dto.IndicatorToMoId)
	// form.Set("indicator_to_mo_fact_id", dto.IndicatorToMoFactId)
	// form.Set("value", dto.Value)
	// form.Set("fact_time", dto.FactTime)
	// form.Set("is_plan", dto.IsPlan)
	// form.Set("auth_user_id", dto.AuthUser)
	// form.Set("comment", dto.Comment)
	// log.Debug().Str("Action", "Request").Str("Body", fmt.Sprintf("%v", form)).Send()

	// req, err := http.NewRequest("POST", cfg.Server.Url, strings.NewReader(form.Encode()))

	// if err != nil {
	// 	log.Fatal().Err(err).Msg("error in construct POST request")
	// }
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatal().Err(err).Msg("error in construct POST request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.Server.Token))

	resp, err := client.Do(req)
	log.Debug().Str("Action", "Response").Int("StatusCode", resp.StatusCode).Send()

	if resp.StatusCode != 200 {
		log.Debug().Str("Response", "Fail").Str("Response", fmt.Sprintf("%v", resp))
	}
}
