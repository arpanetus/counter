package router

import (
	"github.com/arpanetus/counter/pkg/service"
	"log"
	"net/http"
	"strconv"
)

const numBase = 10

type CounterHandler struct {
	svc service.CounterServicer
}

func New(svc service.CounterServicer) *CounterHandler {
	return &CounterHandler{svc}
}

// ServeHTTP handles any request and updates counter returning the sum in a given duration.
func (h *CounterHandler) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	c := h.svc.Count()
	h.svc.Add()

	_, err := resp.Write([]byte(strconv.FormatUint(c, numBase)))
	if err != nil {
		log.Printf("cannot write into response: %v", err)
	}
}

