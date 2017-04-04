package govuegui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(NewRouter())
	defer ts.Close()
	res, err := http.Get(ts.URL + PathPrefix + "/vue.min.js")
	fmt.Println(res.StatusCode)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Did not found vue library")
	}

}
