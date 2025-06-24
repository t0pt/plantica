package events

import (
	"fmt"
	"time"
)

type Event struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
}

type Date struct {
	Year  int        `json:"year,omitempty"`
	Month time.Month `json:"month,omitempty"`
	Day   int        `json:"day,omitempty"`
}

func (d Date) String() string {
	return fmt.Sprintf("%d-%s-%d", d.Day, d.Month.String(), d.Year)
}

var Events = map[Date][]Event{
	Date{
		Day:   1,
		Month: 1,
	}: []Event{
		{Name: "workkkkkkkkkkkkkkk",
			Description: "have to work now on to do app",
			Date:        time.Now()},
		{Name: "study",
			Description: "have to study now on to do app",
			Date:        time.Now()},
	},
	TodayDate(): []Event{
		{Name: "work",
			Description: "have to work now on to do app",
			Date:        time.Now()},
		{Name: "studyyyyy",
			Description: "have to study now on to do app",
			Date:        time.Now()},
	},
	// time.Now(): []Event{},
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
