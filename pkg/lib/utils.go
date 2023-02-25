package lib

import (
	"fmt"
	"strings"
	"time"
)

func TtmInDays(in_seconds int64) float64 {
	return float64(in_seconds - time.Now().Unix()) / (3600.0 * 24.0)
}

func parseMonth(m string) string {
	if len(m) == 1 {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}

func ParserStringDate(date string) (time.Time, int64) {
	var in_seconds int64
	splited := strings.Split(date, "-")
	month := parseMonth(splited[1])
	date = fmt.Sprintf("%s-%s-%s", splited[0], month, splited[2])
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic("could not parse expiration- correct format is yyyy-mm-dd")
	}

	in_seconds = dt.Unix()
	return dt, in_seconds
}
