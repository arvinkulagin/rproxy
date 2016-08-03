package transport

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestRegexpTransportText(t *testing.T) {
	// Testing backend server, that will be target for reverse proxy.
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Lorem ipsum dolor set amet")
	}))
	defer target.Close()

	// Regular expression for replacement substrings in target server response.
	re := regexp.MustCompile("[I,i]psum")

	//
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Error(err)
	}

	// Handler for proxy server with custom transport.
	handler := httputil.NewSingleHostReverseProxy(targetURL)
	handler.Transport = NewRegexpTransport(re, "shmipsum")
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	resp, err := http.Get(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Lorem shmipsum dolor set amet\n" {
		t.Errorf("wrong response body: %s", string(body))
	}
	if resp.ContentLength != 30 {
		t.Errorf("wrong response content length: %d", resp.ContentLength)
	}
}

func TestRegexpTransportJSON(t *testing.T) {
	// Testing backend server, that will be target for reverse proxy.
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"text":"Lorem ipsum dolor set amet"}`)
	}))
	defer target.Close()

	// Regular expression for replacement substrings in target server response.
	re := regexp.MustCompile("[I,i]psum")

	//
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Error(err)
	}

	// Handler for proxy server with custom transport.
	handler := httputil.NewSingleHostReverseProxy(targetURL)
	handler.Transport = NewRegexpTransport(re, "shmipsum")
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	resp, err := http.Get(proxy.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"text":"Lorem ipsum dolor set amet"}`+"\n" {
		t.Errorf("wrong response body: %s", string(body))
	}
	if resp.ContentLength != 38 {
		t.Errorf("wrong response content length: %d", resp.ContentLength)
	}
}

func TestRegexpTransportPost(t *testing.T) {
	// Testing backend server, that will be target for reverse proxy.
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	}))
	defer target.Close()

	// Regular expression for replacement substrings in target server response.
	re := regexp.MustCompile("[I,i]psum")

	//
	targetURL, err := url.Parse(target.URL)
	if err != nil {
		t.Error(err)
	}

	// Handler for proxy server with custom transport.
	handler := httputil.NewSingleHostReverseProxy(targetURL)
	handler.Transport = NewRegexpTransport(re, "shmipsum")
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	resp, err := http.Post(proxy.URL, "text/plain", strings.NewReader("Lorem ipsum dolor set amet\n"))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "Lorem shmipsum dolor set amet\n" {
		t.Errorf("wrong response body: %s", string(body))
	}
	if resp.ContentLength != 30 {
		t.Errorf("wrong response content length: %d", resp.ContentLength)
	}
}
