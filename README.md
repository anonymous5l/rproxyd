# rproxyd

very simple `http` `https` reverse proxy like `nginx` `proxy_pass`

## Build/Usage

```bash
$ go install github.com/anonymous5l/rproxyd
```

default bind port 8080 `url` argument must set

```bash
$ ./rproxyd --bind :8080 --url http://www.baidu.com
```
