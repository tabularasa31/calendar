package entity

import (
	"time"

	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/date_utils"
)

type (
	EventDB struct {
		ID           int
		Title        string
		Desc         string
		UserID       int
		StartTime    time.Time
		EndTime      time.Time
		Notification time.Time
	}
)

func (e *EventDB) Dto() *Event {
	date := date_utils.TimeToString(e.StartTime)

	d := date_utils.TimeToString(e.EndTime)

	n := date_utils.TimeToString(e.Notification)

	event := Event{
		ID:           e.ID,
		Title:        e.Title,
		Desc:         e.Desc,
		UserID:       e.UserID,
		StartTime:    date,
		EndTime:      d,
		Notification: n,
	}

	return &event
}