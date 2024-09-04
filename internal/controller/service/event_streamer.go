package service

import (
	"context"
	dom "payment-sse/internal/domain"
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/env"
	"payment-sse/internal/repo"
	"payment-sse/internal/util"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EventStreamer interface {
	NewStream(e env.Env, ordId uuid.UUID) <-chan domOrd.PaymentEvent
}

type client struct {
	lastUpdate time.Time
	rx         chan<- domOrd.PaymentEvent
}

func newClient(rx chan<- domOrd.PaymentEvent) client {
	lastUpdate := time.Now()
	return client{lastUpdate, rx}
}

type eventStreamer[T any] struct {
	ordRepo repo.OrderRepo[T]
	log     *zap.Logger
	domConf dom.Config

	mu    sync.Mutex
	watch map[uuid.UUID][]client
	cache map[uuid.UUID]struct {
		UserID    uuid.UUID
		CreatedAt time.Time

		Events []domOrd.EventOrder
	}
}

func newCacheEntry(userId uuid.UUID, createdAt time.Time, events []domOrd.EventOrder) struct {
	UserID    uuid.UUID
	CreatedAt time.Time

	Events []domOrd.EventOrder
} {
	return struct {
		UserID    uuid.UUID
		CreatedAt time.Time

		Events []domOrd.EventOrder
	}{UserID: userId, CreatedAt: createdAt, Events: events}
}

func NewEventStreamer[T any](
	ordRepo repo.OrderRepo[T],
	log *zap.Logger,
	domConf dom.Config,
) EventStreamer {
	es := &eventStreamer[T]{
		ordRepo: ordRepo,
		log:     log,
		domConf: domConf,
		mu:      sync.Mutex{},
		watch:   make(map[uuid.UUID][]client),
		cache: make(map[uuid.UUID]struct {
			UserID    uuid.UUID
			CreatedAt time.Time

			Events []domOrd.EventOrder
		}),
	}

	go es.startPolling()

	return es
}

func (s *eventStreamer[T]) startPolling() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		s.removeStale()
		followed := make([]uuid.UUID, 0, len(s.watch))
		for k := range s.watch {
			followed = append(followed, k)
		}
		s.mu.Unlock()

		func() {
			if len(followed) == 0 {
				return
			}

			ctx, canc := context.WithTimeout(context.Background(), 700*time.Millisecond)
			defer canc()

			e := env.NewEnv(ctx, s.log, &s.domConf)
			timeOrdEvents, rerr := s.ordRepo.GetChronOrdersEvents(e, followed)
			if rerr != nil {
				e.LogUnexpectedError(rerr, "polling for updated order events")
				return
			}

			for id, ord := range timeOrdEvents {
				s.update(e, id, &ord)
			}
		}()

	}
}

// Assumes the events are correct and sorted by their creation time
func (s *eventStreamer[T]) update(e env.Env, ordId uuid.UUID, updOrd *struct {
	UserID    uuid.UUID
	CreatedAt time.Time

	Events []domOrd.EventOrder
},
) {
	inOrder := u.Fold(
		updOrd.Events,
		func(acc []domOrd.EventOrder, curr domOrd.EventOrder) []domOrd.EventOrder {
			// Include non-final events going one after the other
			if curr.Status <= domOrd.ConfirmedByMayor && int(curr.Status) == len(acc) {
				acc = append(acc, curr)
				return acc
			}

			// Include final events
			if len(acc) >= 3 {
				acc = append(acc, curr)
				return acc
			}

			return acc
		},
		make([]domOrd.EventOrder, 0, len(updOrd.Events)),
	)

	s.mu.Lock()
	defer s.mu.Unlock()

	newEvents := func() []domOrd.EventOrder {
		oldOrd, ok := s.cache[ordId]
		if !ok {
			return inOrder
		}

		if len(inOrder) <= len(oldOrd.Events) {
			return nil
		}

		start, finish := len(oldOrd.Events), len(inOrder)
		return inOrder[start:finish]
	}()
	if len(newEvents) == 0 {
		return
	}
	e.LogDebug("new events encountered", newEvents)
	s.cache[ordId] = newCacheEntry(updOrd.UserID, updOrd.CreatedAt, inOrder)

	clients, ok := s.watch[ordId]
	if !ok {
		return
	}

	isFin := newEvents[len(newEvents)-1].Final(e.DomConf().PaymentConfirmationIn)
	for i := range clients {
		for _, eve := range newEvents {
			// Should never block since the buffer is set to 5
			s.watch[ordId][i].rx <- domOrd.NewPaymentEvent(ordId, updOrd.UserID, updOrd.CreatedAt, eve)
		}
		// Just disconnect, if it is final we will delete all watching clients immediately
		if isFin {
			close(s.watch[ordId][i].rx)
			continue
		}

		s.watch[ordId][i].lastUpdate = time.Now()
	}

	if isFin {
		delete(s.watch, ordId)
	}
}

// SAFETY: Requires an active lock
func (s *eventStreamer[T]) removeStale() {
	for ordId, clients := range s.watch {
		nonStale := make([]client, 0, len(clients))
		for _, cli := range clients {
			if time.Since(cli.lastUpdate) > s.domConf.InactivityOrderEventTimeout {
				close(cli.rx)
				continue
			}
			nonStale = append(nonStale, cli)
		}
		if len(nonStale) == 0 {
			delete(s.watch, ordId)
			delete(s.cache, ordId)
			continue
		}

		s.watch[ordId] = nonStale
	}
}

func (s *eventStreamer[T]) NewStream(e env.Env, ordId uuid.UUID) <-chan domOrd.PaymentEvent {
	tx := make(chan domOrd.PaymentEvent, 5)

	s.mu.Lock()
	clis, ok := s.watch[ordId]
	if !ok {
		clis = make([]client, 0)
	}
	clis = append(clis, newClient(tx))
	s.watch[ordId] = clis

	ord, ok := s.cache[ordId]
	s.mu.Unlock()

	if ok {
		events := ord.Events
		go func() {
			e.LogDebug("retrieving already fetched", events)
			for _, eve := range events {
				tx <- domOrd.NewPaymentEvent(ordId, ord.UserID, ord.CreatedAt, eve)
			}
		}()
	}

	return tx
}
