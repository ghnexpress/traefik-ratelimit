package utils

import (
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func GetUrl(host string, pathJoined ...string) (*url.URL, error) {
	hostParsed, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	hostParsed.Path = path.Join(pathJoined...)
	return hostParsed, nil
}

func GetIp(req *http.Request) string {
	if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			return strings.TrimSpace(ip[:i])
		}
		return ip
	}
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ra
}
