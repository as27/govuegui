package govuegui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterWithRice(t *testing.T) {
	useRice = true
	routerTest(t)
	useRice = false
	routerTest(t)
}

func routerTest(t *testing.T) {
	ts := httptest.NewServer(NewRouter())
	defer ts.Close()
	testUrls := []string{
		"lib/vue.min.js",
		"lib/pure.min.css",
		"lib/app.js",
	}
	for _, tURL := range testUrls {
		res, err := http.Get(ts.URL + PathPrefix + "/" + tURL)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != 200 {
			t.Errorf("Did not find %s", tURL)
		}
	}
}
