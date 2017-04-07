package govuegui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {

	ts := httptest.NewServer(NewRouter())
	defer ts.Close()
	res, err := http.Get(ts.URL + PathPrefix + "/vue.min.js")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Did not find vue library")
	}
	res, err = http.Get(ts.URL + PathPrefix + "/app.js")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Did not find app.js")
	}

}