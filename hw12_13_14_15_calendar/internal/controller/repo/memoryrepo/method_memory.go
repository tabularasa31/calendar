package memoryrepo

import (
	"context"
	"github.com/google/uuid"
	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"sync"
	"time"
)

// EventRepo -.
type EventRepo struct {
	events map[int]map[int]entity.EventDB
	mu     *sync.RWMutex
}

// New -.
func New() *EventRepo {
	m := sync.RWMutex{}
	events := make(map[int]map[int]entity.EventDB)
	return &EventRepo{
		events: events,
		mu:     &m,
	}
}

// CreateEvent Создать (событие).
func (r *EventRepo) CreateEvent(ctx context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if event time already busy
	userEvents, ok := r.events[eventDB.UserID]
	if !ok {
		r.events[eventDB.UserID] = make(map[int]entity.EventDB)
	} else if !r.isEventTimeBusy(userEvents, *eventDB) {
		return nil, errapp.ErrEventTimeBusy
	}

	// Create unique event ID
	eventDB.ID = int(uuid.New().ID())

	r.events[eventDB.UserID][eventDB.ID] = *eventDB
	res := r.events[eventDB.UserID][eventDB.ID]

	return res.Dto(), nil
}

// UpdateEvent Обновить (ID пользователя, ID события, событие).
func (r *EventRepo) UpdateEvent(ctx context.Context, eventDB *entity.EventDB) (*entity.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userEvents, ok := r.events[eventDB.UserID]
	if !ok {
		return nil, errapp.ErrEventNotFound
	}

	updatedEvent, ok := userEvents[eventDB.ID]
	if !ok {
		return nil, errapp.ErrEventNotFound
	}

	if !r.isEventTimeBusy(userEvents, *eventDB) {
		return nil, errapp.ErrEventTimeBusy
	}

	updatedEvent.Title = eventDB.Title
	updatedEvent.Desc = eventDB.Desc
	updatedEvent.StartTime = eventDB.StartTime
	updatedEvent.EndTime = eventDB.EndTime

	r.events[eventDB.UserID][eventDB.ID] = updatedEvent
	res := r.events[eventDB.UserID][eventDB.ID]
	return res.Dto(), nil
}

// DeleteEvent Удалить (ID события).
func (r *EventRepo) DeleteEvent(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, userEvents := range r.events {
		if _, ok := userEvents[id]; !ok {
			delete(userEvents, id)
			return nil
		}
	}
	return errapp.ErrEventNotFound
}

// GetDailyEvents СписокСобытийНаДень (дата).
// Выводит все события, которые начинаются в заданный день.
//func (r *EventRepo) GetDailyEvents(ctx context.Context, userID int, date time.Time) ([]entity.Event, error) {
//	var events []entity.Event
//
//	r.mu.RLock()
//	defer r.mu.RUnlock()
//
//	for _, userEvents := range r.events {
//		for _, evDB := range userEvents {
//			if evDB.StartTime.Day() == date.Day() {
//				events = append(events, *evDB.Dto())
//			}
//		}
//	}
//	return events, nil
//}

func (r *EventRepo) GetEventsByDates(ctx context.Context, uid int, start time.Time, end time.Time) ([]entity.Event, error) {
	var userEvents []entity.Event

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, userEvent := range r.events[uid] {
		if userEvent.StartTime.After(start) && userEvent.StartTime.Before(end) {
			userEvents = append(userEvents, *userEvent.Dto())
		}
	}

	return userEvents, nil
}

// isEventTimeBusy проверка на занятость в заданное время.
func (r *EventRepo) isEventTimeBusy(userEvents map[int]entity.EventDB, newEvent entity.EventDB) bool {
	newStartTime := newEvent.StartTime
	newEndTime := newEvent.EndTime
	for _, userEvent := range userEvents {
		oldStartTime := userEvent.StartTime
		oldEndTime := userEvent.EndTime
		if (newStartTime.After(oldStartTime) && newStartTime.Before(oldEndTime)) ||
			(newEndTime.After(oldStartTime) && newEndTime.Before(oldEndTime)) {
			return false
		}
	}
	return true
}
