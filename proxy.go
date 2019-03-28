package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

// writeFile read local file and write back
func (p *ReverseProxyHandler) writeFile(request *fasthttp.RequestCtx, file string) {
	f, err := os.Open(file)

	if err != nil {
		request.SetStatusCode(403)
		request.SetBody([]byte(fmt.Sprint("can't open file %s err %s", file, err)))
		return
	}

	request.SetBodyStream(f, -1)
}

func (p *ReverseProxyHandler) listDirectory(request *fasthttp.RequestCtx, dir string) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		request.SetStatusCode(403)
		request.SetBody([]byte(fmt.Sprint("can't open dir %s err %s", dir, err)))
		return
	}

	te := TemplateEntity{}
	te.SetTitle(dir)

	for _, file := range files {
		// filter hidden file
		if file.Name()[0] != '.' {
			te.AppendItem(file, filepath.Join("/", dir, file.Name()))
		}
	}

	te.Sort()

	err = Template.Execute(request.Response.BodyWriter(), te)

	if err != nil {
		request.SetStatusCode(500)
		request.SetBody([]byte(fmt.Sprint("can't renderer template dir %s err %s", dir, err)))
	}
}

func (p *ReverseProxyHandler) HandleIndex(request *fasthttp.RequestCtx, u url.URL, method string) {
	request.SetContentType("text/html; charset=utf-8")

	fi, err := os.Stat(u.Path)

	if os.IsNotExist(err) {
		request.SetStatusCode(404)
		request.SetBody([]byte(fmt.Sprintf("%s not found!", u.Path)))
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		// list directory
		p.listDirectory(request, u.Path)
	case mode.IsRegular():
		request.SetContentType("application/octet-stream")
		p.writeFile(request, u.Path)
		return
	}
}

func (p *ReverseProxyHandler) Handle(request *fasthttp.RequestCtx) {
	u := p.proxyUrl

	spath := string(request.Path())

	// try fix root path
	if spath == "/" {
		u.Path = p.proxyUrl.Path
	} else if filepath.IsAbs(spath) {
		u.Path = spath
	} else {
		u.Path = filepath.Join(p.proxyUrl.Path, spath)
	}

	u.RawQuery = string(request.URI().QueryString())

	smethod := string(request.Method())

	console.Log("%s %s", smethod, u.String())

	if u.Scheme == "" && u.Host == "" && smethod == "GET" {
		p.HandleIndex(request, u, smethod)
		return
	}

	postBody := request.PostBody()

	req, err := http.NewRequest(smethod, u.String(), bytes.NewReader(postBody))

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
