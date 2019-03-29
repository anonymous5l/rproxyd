# rproxyd

very simple `http` `https` reverse proxy like `nginx` `proxy_pass`

also support local static resource mapping

## Build/Usage

```bash
$ go get -u github.com/anonymous5l/rproxyd
```

default bind port 8080 `url` argument must set

```bash
$ rproxyd --bind :8080 --url http://www.baidu.com
```

for static resource mapping

```bash
$ rproxyd --bind :8080 --url ./
```

or

```bash
$ rproxyd --bind :8080 --url /Users/xxxx/Desktop
```
