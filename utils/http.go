package utils

import (
	"net/url"
	"path"
)

func GetUrl(host string, pathJoined ...string) (*url.URL, error) {
	hostParsed, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	hostParsed.Path = path.Join(pathJoined...)
	return hostParsed, nil
}
