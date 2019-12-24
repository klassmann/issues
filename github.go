package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// IssueURL is the URL to github API
const IssueURL = "https://api.github.com/search/issues"

// User information struct
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// Issue information struct
type Issue struct {
	Number    int
	HTMLURL   string
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

// IssueQueryResult response struct from github api
type IssueQueryResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
	NextPage   string
	LastPage   string
}

func (i *IssueQueryResult) addLink(s string) {
	if strings.Contains(s, "next") {
		i.NextPage = s
	} else if strings.Contains(s, "last") {
		i.LastPage = s
	}
}

func httpGet(terms []string) (*http.Response, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssueURL + "?q=" + q)

	if err != nil {
		resp.Body.Close()
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	return resp, nil
}

func tryPaginationLinks(r *http.Response) []string {
	l := r.Header.Get("Link")
	return strings.Split(l, ",")
}

func queryIssues(terms []string) (*IssueQueryResult, error) {
	resp, err := httpGet(terms)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result IssueQueryResult

	// Add pagination links
	links := tryPaginationLinks(resp)
	for i := 0; i < len(links); i++ {
		result.addLink(links[i])
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
