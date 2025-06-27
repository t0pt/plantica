package events

import (
	"database/sql"
	"fmt"
	"time"
)

type EventManager struct {
	DbPath string
	db     *sql.DB
}

type Event struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Time        int    `json:"time,omitempty"`
}

type Date struct {
	Year  int        `json:"year,omitempty"`
	Month time.Month `json:"month,omitempty"`
	Day   int        `json:"day,omitempty"`
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%s-%d", d.Day, d.Month.String(), d.Year)
}

func (d Date) Int() int64 {
	return d.ToTime().Unix()
}

func TodayDate() Date {
	now := time.Now()
	return Date{
		Year:  now.Year(),
		Month: now.Month(),
		Day:   now.Day(),
	}
}

func (d Date) AddDays(n int) Date {
	return TimeToDate(d.ToTime().AddDate(0, 0, n))
}

func (d Date) ToTime() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

func TimeToDate(input time.Time) Date {
	return Date{
		Year:  input.Year(),
		Month: input.Month(),
		Day:   input.Day(),
	}
}
