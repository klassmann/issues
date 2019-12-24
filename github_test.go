package main

import (
	"net/http"
	"testing"
)

func mockResponse() *http.Response {
	r := http.Response{}
	r.Header = http.Header{}
	return &r
}

func TestTryPaginationLinksTwoLinks(t *testing.T) {
	r := mockResponse()
	r.Header.Add("Link", "<https://api.github.com/search/issues?q=windows+label%3Abug+language%3Apython+state%3Aopen&sort=created&order=asc&page=2>; rel=\"next\", <https://api.github.com/search/issues?q=windows+label%3Abug+language%3Apython+state%3Aopen&sort=created&order=asc&page=34>; rel=\"last\"")

	links := tryPaginationLinks(r)

	if len(links) != 2 {
		t.Errorf("links: expected 2 but got %d", len(links))
	}
}
