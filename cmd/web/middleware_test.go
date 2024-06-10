package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {


  rr := httptest.NewRecorder()

  r,err := http.NewRequest(httpMethodGet, "/",  nil)  

  if err != nil {
    t.Fatal(err)
  }


  next := http.HandlerFunc
}
