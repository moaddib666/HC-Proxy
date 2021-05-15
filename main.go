package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var InjectedScriptSRC string
var BackendURL string
var DefaultPort string
var DefaultListenAddress string

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	BackendURL = GetEnv("HC_PROXY_BACKEND_URL", "http://hellcorp.com.ua")
	InjectedScriptSRC = GetEnv("HC_PROXY_INJECTION_SCRIPT_SRC", BackendURL+"/js/label.js")
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
	if requestMethod != "GET" {
		log.Println("Unsupported Method", requestMethod, "in request")
		return nil
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		log.Println("Not html response")
		return nil
	}
	log.Println("Processing", requestURL)
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	if strings.Contains(bodyString, "</body>") {
		log.Println("Response ", requestURL, " contains html body tag")
		bodyString = strings.ReplaceAll(bodyString, "</body>", "<script src='"+InjectedScriptSRC+"'></script></body>")
		bodyBytes = []byte(bodyString)
		bodyContentLength := len(bodyBytes)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		r.ContentLength = int64(bodyContentLength)
		r.Header.Set("Content-Length", strconv.Itoa(bodyContentLength))
	}

	return nil
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-HC-Proxy", "True")
		p.ServeHTTP(w, r)
	}
}
