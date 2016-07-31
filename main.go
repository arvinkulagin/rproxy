package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"

	"github.com/arvinkulagin/rproxy/config"
	"github.com/arvinkulagin/rproxy/transport"
)

var filename = flag.String("config", "", "Config file")

func main() {
	flag.Parse()

	if *filename == "" {
		log.Fatal("Specify config file")
	}

	// Open JSON config file.
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read from file and unmarshal raw JSON.
	conf, err := config.ReadConfig(file)
	if err != nil {
		log.Fatal(err)
	}

	// Try to compile specified regexp pattern.
	re, err := regexp.Compile(conf.Pattern)
	if err != nil {
		log.Fatal(err)
	}

	// Try to parse specified
	targetURL, err := url.Parse(conf.Target)
	if err != nil {
		log.Fatal(err)
	}

	// Get default HTTP proxy handler.
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Use RegexpTransport round-tripper if regexp pattern is specified.
	if conf.Pattern != "" {
		proxy.Transport = transport.NewRegexpTransport(re, conf.Replacement)
	}

	fmt.Println("Listening on " + conf.Address)

	log.Fatal(http.ListenAndServe(conf.Address, proxy))
}
