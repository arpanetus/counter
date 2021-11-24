package main

import (
	"fmt"
	"github.com/arpanetus/counter/pkg/file"
	"github.com/arpanetus/counter/pkg/router"
	"github.com/arpanetus/counter/pkg/service"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestCounterMux(t *testing.T) {
	// Generate some random file.
	f, err := ioutil.TempFile("", "counter*")
	if err != nil {
		t.Fatalf("err occured while testing: %+v", err)
	}

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Printf("cannot remove temp file after test: %+v", err)
		}
	}(f.Name())

	fmt.Println("temp file created")

	path, err := filepath.Abs(f.Name())

	// Add some past value.
	pastCountNano := time.Now().Add(-time.Minute).UnixNano()
	pastCount := strconv.FormatInt(pastCountNano, 10) + "\n"
	_, _ = f.WriteString(pastCount)

	// Pass that guy further.
	fw := file.New(path, Minute, file.FS)
	svc := service.New(fw)
	if err := svc.Parse(); err != nil {
		log.Fatalf("cannot parse times on init: %v", err)
	}

	handler := router.New(svc)

	req := httptest.NewRequest(http.MethodGet, "/count", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s!=nil", err)
	}

	// Check whether answer is 1.
	if dataStr := string(data); dataStr != fmt.Sprintf(router.CountStr, 1) {
		t.Errorf("expected 1 count, got: %s", dataStr)
	}

	// Let's read what has been written into file.
	f.Seek(0, 0)

	b, _ := fw.Read()
	stamps := b.Stamps()
	if l := len(stamps); l != 1 {
		t.Errorf("expected 1 stamp, got %d", l)
	}

	var gotCountNano int64
	for k, _ := range stamps {
		gotCountNano = k
	}

	// Now let's check if gotCountNano is greater than the past one.
	if gotCountNano <= pastCountNano {
		t.Errorf("got num is smaller than past entered number :<")
	}

}
