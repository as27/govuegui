package govuegui

import (
	"fmt"
	"io/ioutil"
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
		t.Error("Did not found vue library")
	}

	res, err = http.Get(ts.URL + PathPrefix + "/app.js")
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}
