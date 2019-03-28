module github.com/anonymous5l/rproxyd

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190313024323-a1f597ede03a
	golang.org/x/net => github.com/golang/net v0.0.0-20190313220215-9f648a60d977
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190312061237-fead79001313
	golang.org/x/text => github.com/golang/text v0.3.1-0.20190306152657-5d731a35f486
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190314010720-1286b2016bb1
)

require (
	github.com/anonymous5l/console v0.0.0-20190221092207-cdcee6db6f29
	github.com/gabriel-vasile/mimetype v0.1.3
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasthttp v1.2.0
)
