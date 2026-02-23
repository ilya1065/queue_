package parser

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"queue/internal/entity"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

type Response struct {
	PageProps struct {
		ScheduleLoadInfo []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			ICalContent string `json:"iCalContent"`
			ICalLink    string `json:"iCalLink"`
		} `json:"scheduleLoadInfo"`
	} `json:"pageProps"`
}

func ICSURL(start, end time.Time, url string) ([]entity.ScheduleItem, error) {
	slog.Info("Запрос к API и парсинг расписания")
	// запрос к api университета
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var body Response
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if len(body.PageProps.ScheduleLoadInfo) == 0 {
		return nil, errors.New("ошибка получения расписания")
	}
	r := strings.NewReader(body.PageProps.ScheduleLoadInfo[0].ICalContent)
	c := gocal.NewParser(r)
	c.Start = &start
	c.End = &end
	err = c.Parse()
	if err != nil {
		return nil, err
	}
	var events []entity.ScheduleItem
	for _, ev := range c.Events {
		events = append(events, entity.ScheduleItem{
			Name:        ev.Summary,
			Description: ev.Description,
			StartDate:   ev.Start,
			EndDate:     ev.End,
			ExternalID:  ev.Uid,
		})
	}
	return events, nil
}
