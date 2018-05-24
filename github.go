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

func queryIssues(terms []string) (*IssueQueryResult, error) {
	resp, err := httpGet(terms)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result IssueQueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
