package server

import (
	"net/http/httputil"

	"github.com/JustinLi007/whatdoing/libs/go/utils"
)

func (s *Server) NewReverseProxy() *httputil.ReverseProxy {
	rewriteFn := func(pr *httputil.ProxyRequest) {
		prefix, pathValue, ok := utils.ParseRequestUrl(pr.In)
		if !ok {
			return
		}

		endpoint, err := s.ServiceMap.GetEndpoint(prefix)
		if err != nil {
			return
		}
		endpoint.Url = endpoint.Url.JoinPath("/")
		if pathValue != "" {
			endpoint.Url = endpoint.Url.JoinPath(pathValue)
		}

		pr.Out.URL = endpoint.Url
		pr.Out.Host = pr.In.Host

		pr.Out.Header["X-Forwarded-For"] = pr.In.Header["X-Forwarded-For"]
		pr.SetXForwarded()
	}

	rp := &httputil.ReverseProxy{
		Rewrite: rewriteFn,
	}

	return rp
}
