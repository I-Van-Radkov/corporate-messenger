package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Factory struct {
	cache sync.Map
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Get(target string) *httputil.ReverseProxy {
	if p, ok := f.cache.Load(target); ok {
		return p.(*httputil.ReverseProxy)
	}

	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	origDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		origDirector(r)

		// Сохраняем оригинальный хост
		r.Host = u.Host
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host

		// КОРОННОЕ ИСПРАВЛЕНИЕ: передаем все заголовки оригинала
		// Особенно важны X-User-ID и X-User-Role
		for key, values := range r.Header {
			for _, value := range values {
				r.Header.Add(key, value)
			}
		}

		// Удаляем заголовки, которые могут мешать
		r.Header.Del("Accept-Encoding")

		// Добавляем заголовки для CORS
		r.Header.Set("Access-Control-Allow-Origin", "*")
		r.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		r.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-User-ID,X-User-Role")
		r.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	f.cache.Store(target, proxy)
	return proxy
}
