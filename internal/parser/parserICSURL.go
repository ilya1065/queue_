package parser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"queue/internal/entity"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

//type Response struct {
//	PageProps struct {
//		ScheduleLoadInfo []struct {
//			ID          int    `json:"id"`
//			Title       string `json:"title"`
//			ICalContent string `json:"iCalContent"`
//			ICalLink    string `json:"iCalLink"`
//		} `json:"scheduleLoadInfo"`
//	} `json:"pageProps"`
//}

func ICSURL(start, end time.Time, url string) ([]entity.ScheduleItem, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Статутс: %s", resp.Status)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	exdatesByUID, err := parseEXDATES(raw)
	if err != nil {
		return nil, err
	}

	f := bytes.NewReader(raw)

	c := gocal.NewParser(f)
	c.Start = &start
	c.End = &end
	err = c.Parse()
	if err != nil {
		return nil, err
	}
	var events []entity.ScheduleItem

	for _, ev := range c.Events {
		if isExcluded(exdatesByUID, ev.Uid, ev.Start) {
			continue
		}
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

func parseEXDATES(raw []byte) (map[string]map[time.Time]struct{}, error) {
	lines := unfoldICSLines(string(raw))

	result := make(map[string]map[time.Time]struct{})

	var uid string
	var insideEvent bool
	var pending []string

	for _, line := range lines {
		switch {
		case line == "BEGIN:VEVENT":
			insideEvent = true
			uid = ""
			pending = pending[:0]

		case line == "END:VEVENT":
			if uid != "" {
				for _, exline := range pending {
					if err := applyEXDATELine(result, uid, exline); err != nil {
						return nil, err
					}
				}
			}
			insideEvent = false

		case insideEvent && strings.HasPrefix(line, "UID:"):
			uid = strings.TrimPrefix(line, "UID:")

		case insideEvent && strings.HasPrefix(line, "EXDATE"):
			pending = append(pending, line)
		}
	}

	return result, nil
}

func isExcluded(exdatesByUID map[string]map[time.Time]struct{}, uid string, start *time.Time) bool {
	if start == nil {
		return false
	}
	exdates, ok := exdatesByUID[uid]
	if !ok {
		return false
	}
	_, exists := exdates[start.UTC()]
	return exists
}

func unfoldICSLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	rawLines := strings.Split(s, "\n")

	var lines []string
	for _, line := range rawLines {
		if len(lines) > 0 && (strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")) {
			lines[len(lines)-1] += strings.TrimLeft(line, " \t")
		} else {
			lines = append(lines, line)
		}
	}
	return lines
}

func extractTZID(meta string) string {
	for _, p := range strings.Split(meta, ";") {
		if strings.HasPrefix(p, "TZID=") {
			return strings.TrimPrefix(p, "TZID=")
		}
	}
	return ""
}

func parseICSTime(v string, loc *time.Location) (time.Time, error) {
	switch {
	case strings.HasSuffix(v, "Z"):
		if t, err := time.Parse("20060102T150405Z", v); err == nil {
			return t, nil
		}
		return time.Parse("20060102T1504Z", v)
	case strings.Contains(v, "T"):
		if t, err := time.ParseInLocation("20060102T150405", v, loc); err == nil {
			return t, nil
		}
		return time.ParseInLocation("20060102T1504", v, loc)
	default:
		return time.ParseInLocation("20060102", v, loc)
	}
}

func applyEXDATELine(result map[string]map[time.Time]struct{}, uid, line string) error {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return nil
	}

	meta := parts[0]
	values := parts[1]

	loc := time.UTC
	if strings.Contains(meta, "TZID=") {
		tzid := extractTZID(meta)
		if l, err := time.LoadLocation(tzid); err == nil {
			loc = l
		}
	}

	if result[uid] == nil {
		result[uid] = make(map[time.Time]struct{})
	}

	for _, v := range strings.Split(values, ",") {
		t, err := parseICSTime(v, loc)
		if err != nil {
			continue
		}
		result[uid][t.UTC()] = struct{}{}
	}

	return nil
}

//func ICSURL(start, end time.Time, url string) ([]entity.ScheduleItem, error) {
//	slog.Info("Запрос к API и парсинг расписания")
//	// запрос к api университета
//	resp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return nil, errors.New(resp.Status)
//	}
//	var body Response
//	err = json.NewDecoder(resp.Body).Decode(&body)
//	if err != nil {
//		return nil, err
//	}
//	if len(body.PageProps.ScheduleLoadInfo) == 0 {
//		return nil, errors.New("ошибка получения расписания")
//	}
//	r := strings.NewReader(body.PageProps.ScheduleLoadInfo[0].ICalContent)
//	c := gocal.NewParser(r)
//	c.Start = &start
//	c.End = &end
//	err = c.Parse()
//	if err != nil {
//		return nil, err
//	}
//	var events []entity.ScheduleItem
//	for _, ev := range c.Events {
//		events = append(events, entity.ScheduleItem{
//			Name:        ev.Summary,
//			Description: ev.Description,
//			StartDate:   ev.Start,
//			EndDate:     ev.End,
//			ExternalID:  ev.Uid,
//		})
//	}
//	return events, nil
//}
