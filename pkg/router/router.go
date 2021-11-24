package router

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/arpanetus/counter/pkg/service"
)

type CounterHandler struct {
	svc service.CounterServicer
}

func New(svc service.CounterServicer) *CounterHandler {
	return &CounterHandler{svc: svc}
}

func HasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get(ContentTypeLiteral)

	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

const (
	CountStr = `{
	"count": %d
}`
	MimeAppJson = "application/json"

	ContentTypeLiteral = "Content-Type"
)

// ServeHTTP handles any request and updates counter returning the sum in a given duration.
func (h *CounterHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if HasContentType(req, MimeAppJson) {
		resp.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if req.Method != http.MethodGet {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.svc.Add()
	c := h.svc.Count()
	if err := h.svc.Write(); err != nil {
		ret := fmt.Sprintf("cannot write into counter service: %+v", err)
		log.Print(ret)
		resp.Write([]byte(ret))
		return
	}

	resp.Header().Set(ContentTypeLiteral, MimeAppJson)
	_, err := resp.Write([]byte(fmt.Sprintf(CountStr, c)))
	if err != nil {
		log.Printf("cannot write into response: %v", err)
	}
}
