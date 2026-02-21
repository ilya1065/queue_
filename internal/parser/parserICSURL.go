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

// todo дропать базу и загружать занова
// !
func ICSURL(url string) ([]entity.ScheduleItem, error) {
	slog.Info("Запрос к API и парсинг расписания")
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
	r := strings.NewReader(body.PageProps.ScheduleLoadInfo[0].ICalContent)
	c := gocal.NewParser(r)
	start := time.Now().AddDate(-1, 0, 0)
	c.Start = &start
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
