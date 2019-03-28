package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/anonymous5l/rproxyd/hack"
	"github.com/valyala/fasthttp"

	"github.com/anonymous5l/console"
	"github.com/urfave/cli"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "rproxy"
	app.Usage = "simple reverse proxy"
	app.Version = "beta"
	app.Author = "anonymous5l"
	return app
}

func StartHttpProxy(c *cli.Context) error {
	rawUrl := c.String("url")

	if rawUrl == "" {
		return fmt.Errorf("url can't be empty")
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	handler := NewReverseProxyHandler(*u)

	cert := c.String("cert")
	key := c.String("key")

	var tlsConfig *tls.Config

	if cert != "" && key != "" {
		// enable tls
		certificate, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return err
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{
				certificate,
			},
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
	}

	l, err := net.Listen(c.String("type"), c.String("bind"))
	if err != nil {
		return err
	}

	console.Ok("start listening on %s", c.String("bind"))

	for {
		var fconn net.Conn

		conn, err := l.Accept()

		if err != nil {
			console.Err("net: accept error! %s", err)
			continue
		}

		iconn := hack.NewIdentityConn(conn)

		t, err := iconn.Identify()
		if err != nil {
			_ = iconn.Close()
			continue
		}

		if t == hack.IdentityHttp {
			//console.Log("net: handle http conn")
			fconn = iconn
		} else if t == hack.IdentityHttps {
			if tlsConfig != nil {
				//console.Log("net: handle https conn")
				fconn = tls.Server(iconn, tlsConfig)
			}
		}

		if fconn == nil {
			_ = iconn.Close()
			continue
		}

		go func(conn net.Conn) {
			_ = fasthttp.ServeConn(conn, handler.Handle)
			//if err != nil {
			//	console.Err("handle: %s", err)
			//}
		}(fconn)
	}
}

func main() {
	app := NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "type",
			Value:  "tcp",
			EnvVar: "RPROXY_TYPE",
			Usage:  "proxy listen type include `tcp` `tcp4` `tcp6` `unix` `udp` default `tcp`",
		},
		cli.StringFlag{
			Name:   "bind",
			Value:  ":8080",
			EnvVar: "RPROXY_BIND",
			Usage:  "proxy listen address",
		},
		cli.StringFlag{
			Name:   "url",
			Value:  "",
			EnvVar: "RPROXY_URL",
			Usage:  "proxy pass url",
		},
		cli.StringFlag{
			Name:   "cert",
			Value:  "",
			EnvVar: "RPROXY_TLS_CERT",
			Usage:  "enable tls proxy",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "",
			EnvVar: "RPROXY_TLS_KEY",
			Usage:  "enable tls proxy private key",
		},
	}

	app.Action = StartHttpProxy

	err := app.Run(os.Args)

	if err != nil {
		console.Err("cli: %s", err)
	}
}
