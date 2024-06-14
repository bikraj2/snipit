package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"snipit.bikraj.net/internal/assert"
)

// func TestPing(t *testing.T) {
//
//   rr := httptest.NewRecorder()
//
//   r,err := http.NewRequest(http.MethodGet, "/" , nil)
//
//   if err!=nil {
//     t.Fatal(err)
//   }
//   ping(rr, r )
//
//   rs :=rr.Result()
//
//
//   assert.Equal(t,  rs.StatusCode,  http.StatusOK)
//
//   defer rs.Body.Close()
//
//   body,err := io.ReadAll(rs.Body)
//
//   if err != nil {
//     t.Fatal(err)
//   }
//   bytes.TrimSpace(body)
//   assert.Equal(t, string(body),"OK")
//
// }
//

func TestPing(t *testing.T) {

	// create a new instanc of application struct

	app := application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	ts := httptest.NewTLSServer(app.routes())

	defer ts.Close()

	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID", urlPath: "/snippet/view/1", wantCode: http.StatusOK, wantBody: "An old silent pond..."}, {
			name: "Non-existent ID", urlPath: "/snippet/view/2", wantCode: http.StatusNotFound}, {
			name: "Negative ID", urlPath: "/snippet/view/-1", wantCode: http.StatusNotFound}, {
			name: "Decimal ID", urlPath: "/snippet/view/1.23", wantCode: http.StatusNotFound}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
