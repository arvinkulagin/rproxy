package transport

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strconv"
)

// RegexpTransport replace content of html and plain text responces
// with specified regular expression and replacement string.
type RegexpTransport struct {
	re          *regexp.Regexp
	replacement []byte
}

func NewRegexpTransport(re *regexp.Regexp, replacement string) http.RoundTripper {
	return RegexpTransport{re, []byte(replacement)}
}

func (r RegexpTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	// Get response from target server with default transport.
	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	// Replace patterns only for "OK" responses.
	if response.StatusCode != http.StatusOK {
		return response, nil
	}

	// There is nothing to replace.
	if response.ContentLength == 0 {
		return response, nil
	}

	// Replace only text content.
	mediatype, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return response, nil
	}
	if mediatype != "text/html" && mediatype != "text/plain" {
		return response, nil
	}

	// Read original content from target service.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := response.Body.Close(); err != nil {
		log.Println(err)
	}

	// Replace all regexp matches in original content.
	content := r.re.ReplaceAll(body, r.replacement)

	// Count lengh of content with replacements.
	length := len(content)

	// Set new body and content length for it.
	response.Body = ioutil.NopCloser(bytes.NewReader(content))
	response.ContentLength = int64(length)
	response.Header.Set("Content-Length", strconv.Itoa(length))

	return response, nil
}
