package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/zbronya/free-chat-to-api/httpclient"
	"github.com/zbronya/free-chat-to-api/utils"
	"io"
	"net/url"
	"strings"
)

var client = httpclient.NewReqClient()

var headersFilterList = map[string]bool{
	"x-real-ip": true, "x-forwarded-for": true, "x-forwarded-proto": true,
	"x-forwarded-port": true, "x-forwarded-host": true, "x-forwarded-server": true,
	"cf-warp-tag-id": true, "cf-visitor": true, "cf-ray": true,
	"cf-connecting-ip": true, "cf-ipcountry": true, "cdn-loop": true,
	"remote-host": true, "x-frame-options": true, "x-xss-protection": true,
	"x-content-type-options": true, "content-security-policy": true,
	"host": true, "cookie": true, "connection": true,
	"content-length": true, "content-encoding": true,
	"x-middleware-prefetch": true, "x-nextjs-data": true, "purpose": true,
	"x-forwarded-uri": true, "x-forwarded-path": true,
	"x-forwarded-method": true, "x-forwarded-protocol": true,
	"x-forwarded-scheme": true, "cf-request-id": true,
	"cf-worker": true, "cf-access-client-id": true,
	"cf-access-client-device-type": true, "cf-access-client-device-model": true,
	"cf-access-client-device-name": true, "cf-access-client-device-brand": true,
}

func ReverseProxy(c *gin.Context) {
	var targetURL, _ = url.Parse("https://chat.openai.com")
	targetURL.Path = c.Request.URL.Path
	targetURL.RawQuery = c.Request.URL.RawQuery

	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if _, found := headersFilterList[strings.ToLower(k)]; !found {
			headers[k] = v[0]
		}
	}

	headers["Host"] = targetURL.Host

	body, _ := io.ReadAll(c.Request.Body)
	resp, err := client.Post(targetURL.String(), headers, body)

	if err != nil {
		utils.ErrorResp(c, 500, "fail to proxy", err)
		return
	}

	for k, v := range resp.Header {
		c.Writer.Header()[k] = v
	}

	c.Writer.WriteHeader(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)

}
