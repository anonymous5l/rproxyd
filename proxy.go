package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/anonymous5l/console"

	"github.com/valyala/fasthttp"
)

type ReverseProxyHandler struct {
	proxyUrl     url.URL
	customHeader []string
}

func NewReverseProxyHandler(proxyUrl url.URL, ch []string) *ReverseProxyHandler {
	return &ReverseProxyHandler{
		proxyUrl:     proxyUrl,
		customHeader: ch,
	}
}

func (p *ReverseProxyHandler) Handle(request *fasthttp.RequestCtx) {
	u := p.proxyUrl
	u.Path = filepath.Join(p.proxyUrl.Path, string(request.Path()))
	u.RawQuery = string(request.URI().QueryString())

	postBody := request.PostBody()

	console.Log("%s %s", request.Method(), u.String())

	req, err := http.NewRequest(string(request.Method()), u.String(), bytes.NewReader(postBody))

	if err != nil {
		request.SetStatusCode(500)
		request.SetBody([]byte(fmt.Sprintf("Internal Error: %s", err)))
		return
	}

	request.Request.Header.VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		// key fix
		if sk == "Host" {
			sv = u.Host
		} else if sk == "Origin" || sk == "Referer" {
			if ou, err := url.Parse(sv); err == nil {
				ou.Host = u.Host
				sv = ou.String()
			}
		}
		req.Header.Set(sk, sv)
		//console.Log("request header %s: %s", sk, sv)
	})

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		request.SetStatusCode(500)
		request.SetBody([]byte(fmt.Sprintf("Internal Error: %s", err)))
		return
	}

	for key, value := range resp.Header {
		for _, v := range value {
			request.Response.Header.Set(key, v)
			//console.Log("response header %s %s", key, v)
		}
	}

	// replace with customHeader
	for _, v := range p.customHeader {
		if idx := strings.Index(v, ":"); idx > 0 {
			request.Response.Header.Set(v[:idx], v[idx+1:])
		}
	}

	request.SetStatusCode(resp.StatusCode)
	request.SetBodyStream(resp.Body, -1)
}
