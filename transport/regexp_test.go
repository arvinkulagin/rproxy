package transport

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"testing"
)

func TestRegexpTransport(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Lorem ipsum dolor set amet")
	}))
	defer target.Close()

	re := regexp.MustCompile("[I,i]psum")
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Error(err)
	}
	handler := httputil.NewSingleHostReverseProxy(targetURL)
	handler.Transport = NewRegexpTransport(re, "shmipsum")
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	resp, err := http.Get(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status is not OK: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if string(body) != "Lorem shmipsum dolor set amet\n" {
		t.Errorf("wrong response body: %s", string(body))
	}

	if resp.ContentLength != 30 {
		t.Errorf("wrong response content length: %d", resp.ContentLength)
	}
}
