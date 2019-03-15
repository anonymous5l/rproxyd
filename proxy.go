package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/anonymous5l/console"

	"github.com/valyala/fasthttp"
)

type ReverseProxyHandler struct {
	proxyUrl url.URL
}

func NewReverseProxyHandler(proxyUrl url.URL) *ReverseProxyHandler {
	return &ReverseProxyHandler{
		proxyUrl: proxyUrl,
	}
}

func (p *ReverseProxyHandler) Handle(request *fasthttp.RequestCtx) {
	u := p.proxyUrl
	u.Path = filepath.Join(p.proxyUrl.Path, string(request.Path()))

	console.Log("request %s", u.String())

	req, err := http.NewRequest(string(request.Method()), u.String(), bytes.NewReader(request.PostBody()))

	if err != nil {
		request.SetStatusCode(500)
		request.SetBody([]byte(fmt.Sprintf("Internal Error: %s", err)))
		return
	}

	request.Request.Header.VisitAll(func(k, v []byte) {
		sk, sv := string(k), string(v)
		req.Header.Set(sk, sv)
		console.Log("request header %s: %s", sk, sv)
	})

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		request.SetStatusCode(500)
		request.SetBody([]byte(fmt.Sprintf("Internal Error: %s", err)))
		return
	}

	defer resp.Body.Close()

	for key, value := range resp.Header {
		for _, v := range value {
			request.Response.Header.Add(key, v)
			console.Log("response header %s %s", key, v)
		}
	}

	request.SetStatusCode(resp.StatusCode)
	//request.SetBodyStream(resp.Body, -1)
	_, err = io.Copy(request.Response.BodyWriter(), resp.Body)
	if err != nil {
		console.Err("io copy: %s", err)
	}
}
