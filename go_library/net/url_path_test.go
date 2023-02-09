package net

import (
	"net/url"
	"testing"
)

func TestUrlPathJoin(t *testing.T) {
	fqdn := "https://www.hellosfsfs.com:9999/who/am/i"
	uri, err := url.Parse(fqdn)
	if err != nil {
		panic(err)
	}
	uri = &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   "/me/too",
	}
	//assert.Equal(t, uri, "https://www.hellosfsfs.com:9999/me/too")
}

func TestUri(t *testing.T) {
	fqdn := "https://www.hellosfsfs.com:9999/who/am/i"
	refer := "hello?hkk=3"
	uri, err := url.Parse(fqdn)
	if err != nil {
		panic(err)
	}
	uri = &url.URL{
		Scheme: uri.Scheme,
		Host:   uri.Host,
		Path:   refer,
	}
	println("uri=" + uri.String())
	//assert.Equal(t, uri, "https://www.hellosfsfs.com:9999/me/too")
}
