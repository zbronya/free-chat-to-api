package httpclient

import "net/http"

type HttpClient interface {
	Get(url string, headers map[string]string) (*http.Response, error)
	Post(url string, headers map[string]string, body []byte) (*http.Response, error)
}
