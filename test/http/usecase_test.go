package http

import (
	"net/http"
	"payment-sse/internal/config"
	domOrd "payment-sse/internal/domain/order"
	"payment-sse/internal/protocol/http/dto/request"
	u "payment-sse/internal/util"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tmaxmax/go-sse"
)

type UsecaseSuite struct {
	conf *config.Config
	t    *Testing
	list *ListOrdersHelper
	pay  *PaymentHelper
}

func newUsecaseSuits(
	conf *config.Config,
	t *Testing,
	list *ListOrdersHelper,
	pay *PaymentHelper,
) UsecaseSuite {
	return UsecaseSuite{conf, t, list, pay}
}

func timeFromOrderStatus(status domOrd.OrderStatus) time.Time {
	base := time.Date(2007, 9, 30, 12, 0, 0, 0, time.Local)

	switch status {
	case domOrd.CoolOrderCreated:
		return base.Add(-60 * time.Second)
	case domOrd.SbuVerificationPending:
		return base.Add(-50 * time.Second)
	case domOrd.ConfirmedByMayor:
		return base.Add(-40 * time.Second)
	case domOrd.Chinazes:
		return base
	default:
		return base.Add(5 * time.Second)
	}
}

func (s *UsecaseSuite) doubleSucc(
	ordId, userId uuid.UUID,
	status domOrd.OrderStatus,
	createdAt time.Time,
) {
	req := &request.PaymentWebhook{
		OrderID:   ordId,
		EventID:   uuid.New(),
		UserID:    userId,
		Status:    status.String(),
		CreatedAt: createdAt,
		UpdatedAt: timeFromOrderStatus(status),
	}
	s.pay.SendPayment(req, http.StatusOK)
	s.pay.SendPayment(req, http.StatusConflict)
}

func (s *UsecaseSuite) doubleFail(
	ordId, userId uuid.UUID,
	status domOrd.OrderStatus,
	expCode int,
	createdAt time.Time,
) {
	req := &request.PaymentWebhook{
		OrderID:   ordId,
		EventID:   uuid.New(),
		UserID:    userId,
		Status:    status.String(),
		CreatedAt: createdAt,
		UpdatedAt: timeFromOrderStatus(status),
	}
	s.pay.SendPayment(req, expCode)
	s.pay.SendPayment(req, expCode)
}

func (s *UsecaseSuite) allInOrder(
	wg *sync.WaitGroup,
	ordId, userId uuid.UUID,
	statuses ...domOrd.OrderStatus,
) sse.EventCallbackRemover {
	rx, canc := s.t.ListenStream(ordId)

	go func() {
		i := 0
		for eve := range rx {
			s.t.Equal(ordId, eve.OrderID)
			s.t.Equal(userId, eve.UserID)
			s.t.Equal(statuses[i], eve.OrderStatus)
			s.t.Equal(timeFromOrderStatus(statuses[i]), eve.UpdatedAt)

			i++
			wg.Done()
		}
	}()

	return canc
}

func (s *UsecaseSuite) RunUsecase() {
	{
		orders := s.list.ListOrders("is_final=true")
		s.t.Equal([]domOrd.Order{}, orders)
	}

	s.inOrderComplete(uuid.New(), uuid.New())
	s.nonFinalsOutOfOrderComplete()

	s.tryToFailAfterChinazes()
	s.tryToChangeMyMindAfterChinazes()

	s.getMoneyBackBeforeChinazes()

	s.tryToGetMoneyBackAfterPaymentTimeoutRanOut()
	s.tryToDoFinalActionsAfterFinal()
	s.tryToDoFinalActionsAfterChinazesRanOut()

	time.Sleep(2 * time.Second)

	{
		orders := s.list.ListOrders("is_final=true")
		s.t.NotEmpty(orders)
		for _, ord := range orders {
			s.t.Equal(true, ord.IsFinal)
		}
	}
	{
		orders := s.list.ListOrders("status=chinazes")
		s.t.NotEmpty(orders)
		for _, ord := range orders {
			s.t.Equal(domOrd.Chinazes, ord.Status)
		}
	}
	{
		orders := s.list.ListOrders("status=chinazes,give_my_money_back")
		s.t.NotEmpty(orders)
		foundMoneyBack := false
		for _, ord := range orders {
			s.t.Contains([]domOrd.OrderStatus{domOrd.Chinazes, domOrd.GiveMyMoneyBack}, ord.Status)
			if ord.Status == domOrd.GiveMyMoneyBack {
				foundMoneyBack = true
			}
		}
		s.t.Equal(true, foundMoneyBack)
	}
	{
		orders := s.list.ListOrders("is_final=true&limit=5")
		s.t.Len(orders, 5)
	}
}

var chinazesCompletion []domOrd.OrderStatus = []domOrd.OrderStatus{
	domOrd.CoolOrderCreated,
	domOrd.SbuVerificationPending,
	domOrd.ConfirmedByMayor,
	domOrd.Chinazes,
}

func (s *UsecaseSuite) inOrderComplete(ordId, userId uuid.UUID) {
	var wg sync.WaitGroup
	wg.Add(len(chinazesCompletion))
	canc := s.allInOrder(&wg, ordId, userId, chinazesCompletion...)
	defer canc()

	now := time.Now()
	s.doubleSucc(ordId, userId, domOrd.CoolOrderCreated, now)
	s.doubleSucc(ordId, userId, domOrd.SbuVerificationPending, now)
	s.doubleSucc(ordId, userId, domOrd.ConfirmedByMayor, now)
	s.doubleSucc(ordId, userId, domOrd.Chinazes, now)
	wg.Wait()
}

func (s *UsecaseSuite) nonFinalsOutOfOrderComplete() {
	nonFin := u.Permut(
		[]domOrd.OrderStatus{
			domOrd.CoolOrderCreated,
			domOrd.SbuVerificationPending,
			domOrd.ConfirmedByMayor,
		},
	)

	for _, nf := range nonFin {
		var wg sync.WaitGroup
		wg.Add(len(chinazesCompletion))
		now := time.Now()
		ordId := uuid.New()
		userId := uuid.New()

		statuses := slices.Concat(nf, []domOrd.OrderStatus{domOrd.Chinazes})

		canc := s.allInOrder(&wg, ordId, userId, chinazesCompletion...)
		for _, status := range statuses {
			s.doubleSucc(ordId, userId, status, now)
		}
		wg.Wait()
		canc()
	}
}

func (s *UsecaseSuite) tryToFailAfterChinazes() {
	now := time.Now()
	ordId := uuid.New()
	userId := uuid.New()

	s.inOrderComplete(ordId, userId)
	s.doubleFail(ordId, userId, domOrd.Failed, http.StatusGone, now)
}

func (s *UsecaseSuite) tryToChangeMyMindAfterChinazes() {
	now := time.Now()
	ordId := uuid.New()
	userId := uuid.New()

	s.inOrderComplete(ordId, userId)
	s.doubleFail(ordId, userId, domOrd.ChangedMyMind, http.StatusGone, now)
}

func (s *UsecaseSuite) getMoneyBackBeforeChinazes() {
	now := time.Now()
	ordId := uuid.New()
	userId := uuid.New()

	s.doubleSucc(ordId, userId, domOrd.CoolOrderCreated, now)
	s.doubleSucc(ordId, userId, domOrd.SbuVerificationPending, now)
	s.doubleSucc(ordId, userId, domOrd.ConfirmedByMayor, now)
	s.doubleSucc(ordId, userId, domOrd.GiveMyMoneyBack, now)
	s.doubleSucc(ordId, userId, domOrd.Chinazes, now)
}

func (s *UsecaseSuite) tryToGetMoneyBackAfterPaymentTimeoutRanOut() {
	var wg sync.WaitGroup
	wg.Add(len(chinazesCompletion))
	now := time.Now()
	ordId := uuid.New()
	userId := uuid.New()

	canc := s.allInOrder(&wg, ordId, userId, chinazesCompletion...)
	s.inOrderComplete(ordId, userId)
	wg.Wait()
	canc()
	time.Sleep(time.Second)
	s.doubleFail(ordId, userId, domOrd.GiveMyMoneyBack, http.StatusGone, now)
}

func (s *UsecaseSuite) tryToDoFinalActionsAfterFinal() {
	fin := u.Permut(
		[]domOrd.OrderStatus{
			domOrd.ChangedMyMind,
			domOrd.Failed,
		},
	)
	for _, f := range fin {
		now := time.Now()
		ordId := uuid.New()
		userId := uuid.New()

		s.doubleSucc(ordId, userId, domOrd.CoolOrderCreated, now)
		s.doubleSucc(ordId, userId, domOrd.SbuVerificationPending, now)
		s.doubleSucc(ordId, userId, domOrd.ConfirmedByMayor, now)

		s.doubleSucc(ordId, userId, f[0], now)
		for _, status := range f[1:] {
			s.doubleFail(ordId, userId, status, http.StatusGone, now)
		}
	}
}

func (s *UsecaseSuite) tryToDoFinalActionsAfterChinazesRanOut() {
	now := time.Now()
	ordId := uuid.New()
	userId := uuid.New()
	s.inOrderComplete(ordId, userId)

	time.Sleep(time.Second)
	for _, status := range []domOrd.OrderStatus{
		domOrd.ChangedMyMind,
		domOrd.Failed,
		domOrd.GiveMyMoneyBack,
	} {
		s.doubleFail(ordId, userId, status, http.StatusGone, now)
	}
}
