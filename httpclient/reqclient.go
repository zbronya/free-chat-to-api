package httpclient

import (
	"github.com/imroc/req/v3"
	"github.com/zbronya/free-chat-to-api/config"
	"net/http"
)

type ReqClient struct {
	client *req.Client
}

func NewReqClient() *ReqClient {
	return &ReqClient{
		client: req.C().
			ImpersonateChrome().SetProxyURL(config.ProxyUrl),
	}
}

func (r *ReqClient) Get(url string, headers map[string]string) (*http.Response, error) {
	resp, err := r.client.R().
		SetHeaders(headers).
		Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}

func (r *ReqClient) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	resp, err := r.client.R().
		SetHeaders(headers).
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	return resp.Response, nil
}
