package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var BackendURL string
var DefaultPort string
var DefaultListenAddress string

const ElementToInject = "<div style=\"border: 1px solid black; border-radius: 6em; background: black; color: white; padding: 1em; font-weight: bold; width: fit-content; box-shadow: rgb(0, 0, 0) 0.1em 0.1em 0.5em; cursor: pointer; position: fixed; left: 1em; bottom: 1em; font-size: 0.6em;z-index: 999;\">Powered by HellCorp Proxy</div>"

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	BackendURL = GetEnv("HC_PROXY_BACKEND_URL", "http://hellcorp.com.ua")
	DefaultPort = GetEnv("HC_PROXY_BIND_PORT", "8980")
	DefaultListenAddress = GetEnv("HC_PROXY_BIND_IP", "0.0.0.0")
	remote, err := url.Parse(BackendURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ModifyResponse = addCustomHeader
	http.HandleFunc("/", handler(proxy))
	log.Println("Starting on http://" + DefaultListenAddress + ":" + DefaultPort)
	err = http.ListenAndServe(DefaultListenAddress+":"+DefaultPort, nil)

	if err != nil {
		panic(err)
	}

}

func addCustomHeader(r *http.Response) error {

	requestMethod := r.Request.Method
	requestURL := r.Request.URL
	contentEncoding := r.Header.Get("Content-Encoding")

	if requestMethod != "GET" {
		log.Println("Unsupported Method", requestMethod, "in request")
		return nil
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		log.Println("Not html response")
		return nil
	}

	log.Println("Processing", requestURL, r.Header)

	if len(contentEncoding) > 1 {
		log.Println("Content Encoded", contentEncoding)
	}

	var reader io.ReadCloser
	switch contentEncoding {
	case "gzip":
		reader, _ = gzip.NewReader(r.Body)
	default:
		reader = r.Body
	}

	bodyBytes, _ := ioutil.ReadAll(reader)
	bodyString := string(bodyBytes[:])
	bodyContentLength := len(bodyBytes)

	if strings.Contains(bodyString, "</body>") {
		log.Println("Response ", requestURL, " contains html body tag")
		bodyString = strings.ReplaceAll(bodyString, "</body>", ElementToInject+"</body>")
		bodyBytes = []byte(bodyString)
		bodyContentLength = len(bodyBytes)
	}

	if contentEncoding == "gzip" {
		bodyBytes = gzipFast(&bodyBytes)
		bodyContentLength = len(bodyBytes)
	}

	writer := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	r.Body = writer
	r.ContentLength = int64(bodyContentLength)
	r.Header.Set("Content-Length", strconv.Itoa(bodyContentLength))

	return nil
}

func gzipFast(a *[]byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(*a); err != nil {
		gz.Close()
		panic(err)
	}
	gz.Close()
	return b.Bytes()
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Accept-Language", "")
		w.Header().Set("X-HC-Proxy", "True")
		p.ServeHTTP(w, r)
	}
}
