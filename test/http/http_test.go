package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment-sse/internal/config"
	"payment-sse/internal/controller"
	"payment-sse/internal/controller/service"
	prot "payment-sse/internal/protocol/http"
	"payment-sse/internal/protocol/http/dto/response"
	pgOrd "payment-sse/internal/repo/postgres/order"
	postgres "payment-sse/internal/repo/postgres/provider"
	pgTx "payment-sse/internal/repo/postgres/transaction"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/tmaxmax/go-sse"
	"go.uber.org/zap"
)

func TestHTTP(t *testing.T) {
	conf, err := config.NewConfig()
	if err != nil {
		t.Fatalf("Failed to source config: %v", err)
	}
	pg, err := postgres.NewPostgresProvider(conf.PG)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := pg.MigrateUp(); err != nil {
		t.Fatalf("Failed to migrate postgres: %v", err)
	}
	log := zap.Must(zap.NewDevelopment())

	ordRepo := pgOrd.NewOrderPostgres(pg.Conn())
	txRepo := pgTx.NewPGTransaction(pg.Conn())
	streamEve := service.NewEventStreamer(
		ordRepo,
		log,
		conf.Dom,
	)

	ordCont := controller.NewOrderController(txRepo, ordRepo, streamEve)

	prot, err := prot.NewHttp(
		conf.Http,
		conf.Dom,
		log,
		ordCont,
	)
	if err != nil {
		t.Fatalf("%v", err)
	}

	go func() {
		prot.Listen()
	}()

	waitForServer(conf.Http.Addr(), t)

	httpCli := http.Client{Timeout: 5 * time.Second}
	test := newTesting(t, &httpCli, conf.Http.Addr(), sse.DefaultClient)
	list := newListOrdersHelper(&conf, test)
	pay := newPaymentHelper(&conf, test)

	useCase := newUsecaseSuits(&conf, test, list, pay)
	useCase.RunUsecase()
}

func waitForServer(baseUrl string, t *testing.T) {
	start := time.Now()
	for {
		resp, err := http.Get(fmt.Sprintf("http://%s/thereisnothing", baseUrl))
		if err == nil && resp.StatusCode == http.StatusNotFound {
			return
		}
		if time.Since(start) > 5*time.Second {
			t.Fatalf("Server did not start within the expected time")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

type Testing struct {
	t       *testing.T
	h       *http.Client
	baseUrl string
	sse     *sse.Client
}

func newTesting(t *testing.T, h *http.Client, baseUrl string, sse *sse.Client) *Testing {
	return &Testing{t, h, baseUrl, sse}
}

func (t *Testing) Nil(object any, msgs ...any) {
	assert.Nil(t.t, object, msgs...)
}

func (t *Testing) Empty(object any, msgs ...any) {
	assert.Empty(t.t, object, msgs...)
}

func (t *Testing) Equal(expected, actual any, msgs ...any) {
	assert.Equal(t.t, expected, actual, msgs...)
}

func (t *Testing) Contains(s, val any, msgs ...any) {
	assert.Contains(t.t, s, val, msgs...)
}

func (t *Testing) NotEmpty(expected any, msg ...any) {
	assert.NotEmpty(t.t, expected, msg)
}

func (t *Testing) Len(obj any, l int, msgs ...any) {
	assert.Len(t.t, obj, l, msgs...)
}

func (t *Testing) Get(path string, exp any, expCode int) {
	resp, err := t.h.Get(path)
	t.Nil(err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	t.Nil(err)

	err = json.Unmarshal(bodyBytes, &exp)
	t.Nil(err)

	if expCode != resp.StatusCode {
		t.t.Errorf(
			"Expected status code: %v, but got: %v and body %v",
			expCode,
			resp.StatusCode,
			string(bodyBytes),
		)
	}
}

func (t *Testing) Post(path string, req any, expCode int) {
	j, err := json.Marshal(req)
	t.Nil(err)

	resp, err := t.h.Post(path, "application/json", bytes.NewBuffer(j))
	t.Nil(err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	t.Nil(err)

	if expCode != resp.StatusCode {
		t.t.Errorf(
			"Expected status code: %v, but got: %v and body %v",
			expCode,
			resp.StatusCode,
			string(bodyBytes),
		)
	}
}

func (t *Testing) ListenStream(
	ordId uuid.UUID,
) (<-chan response.StreamEvent, sse.EventCallbackRemover) {
	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		fmt.Sprintf("http://%s/orders/%s/events", t.baseUrl, ordId),
		http.NoBody,
	)
	t.Nil(err)

	conn := t.sse.NewConnection(req)

	tx := make(chan response.StreamEvent, 5)
	canc := conn.SubscribeEvent("event_order", func(e sse.Event) {
		var event response.StreamEvent
		err = json.Unmarshal([]byte(e.Data), &event)
		t.Nil(err)

		tx <- event
	})

	go func() {
		t.Nil(conn.Connect())
	}()

	return tx, canc
}
