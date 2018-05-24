// Issues helps to search and follow repositories issues list from command line
// You can configure .issuescmdrc to keep saved repositories and predefined searches
// Inspired by https://github.com/adonovan/gopl.io/tree/master/ch4/issues
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

const commandLineHelp = `
Issues help:

add	Adds a new query to your local configuration. You need to give a name before the query,
	$ issues add [query-name] [query parameters]

list	Query issues from a [query-name] saved in local configuration
	$ issues list [query-name]

query	Queries and shows a new search on github.
	- issues query [query-parameters]
	$ issues query repo:golang/go is:open memory

queries Show the list of queries in your configuration
	$ issues queries

help	Shows this help message
`

func addCommand() {

}

func listCommand() {

}

func queryCommand(args []string) {
	result, err := queryIssues(args)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)

	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}

func queriesCommand() {

}

func parseArguments(args []string) {

	if len(args) < 2 || args[1] == "help" {
		fmt.Println(commandLineHelp)
	} else {
		var commandArg = args[1]
		switch commandArg {
		case "add":
			addCommand()
		case "list":
			listCommand()
		case "query":
			queryCommand(args[2:])
		case "queries":
			queriesCommand()
		}
	}
}

func main() {
	parseArguments(os.Args)
}
