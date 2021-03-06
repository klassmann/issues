// Issues helps to search and follow repositories issues list from command line
// You can configure .issuescmdrc to keep saved repositories and predefined searches
// Inspired by https://github.com/adonovan/gopl.io/tree/master/ch4/issues
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	helpAdd        string = "add"
	helpAddDetail  string = "Adds a new query to your local configuration. You need to give a name to the query:"
	helpAddExample string = `$ issues add [query-name] [query parameters]
	$ issues add golang repo:golang/go is:open memory`
	helpDelete        string = "delete"
	helpDeleteDetail  string = "Deletes a query from your local configuration."
	helpDeleteExample string = `$ issues delete [query-name]
	$ issues delete golang`
	helpList         string = "list"
	helpListDetail   string = "Query issues from a [query-name] saved in local configuration:"
	helpListExample  string = `$ issues list [query-name]`
	helpQuery        string = "query"
	helpQueryDetail  string = "Queries and shows a new search on github:"
	helpQueryExample string = `$ issues query [query-parameters]
	$ issues query repo:golang/go is:open memory`
	helpQueries        string = "queries"
	helpQueriesDetail  string = "Show the list of queries in your configuration."
	helpQueriesExample string = "$ issues queries"
	helpHelp           string = "help"
	helpHelpDetail     string = "Shows this help message."
	helpHelpExample    string = ""
)

var (
	termCyan     = TermFormat{FgCyan, AttrBold}
	termBlue     = TermFormat{FgBlue, AttrBold}
	termGreen    = TermFormat{FgGreen, AttrBold}
	termYellow   = TermFormat{FgYellow, AttrReset}
	termGreenDim = TermFormat{FgGreen, AttrDim}
	termRed      = TermFormat{FgRed, AttrBold}
)

type PrintResult struct {
	fromCache bool
	result    *IssueQueryResult
}

func (p *PrintResult) print() {

	if p.fromCache {
		termCyan.Printf("%d issues loaded from cache:\n", p.result.TotalCount)
	} else {
		termCyan.Printf("%d issues loaded from github:\n", p.result.TotalCount)
	}

	for _, item := range p.result.Items {

		number := termBlue.Quote(fmt.Sprintf("#%-5d", item.Number))
		login := termGreenDim.Quote(fmt.Sprintf("%9.9s", item.User.Login))
		title := termYellow.Quote(fmt.Sprintf("%.80s", item.Title))

		fmt.Printf("%s %s %s\n", number, login, title)
	}
}

func addCommand(args []string, conf *Configuration) {

	if len(args) > 1 {
		queryName := args[0]
		query := args[1:]
		conf.Queries[queryName] = strings.Join(query, " ")
		err := saveConfiguration(conf)

		if err != nil {
			termRed.Printf("add command: %v\n", err)
		} else {
			termGreen.Printf("Query %s saved.\n", queryName)
		}
	} else {
		printHelp()
	}
}

func deleteCommand(args []string, conf *Configuration) {
	if len(args) != 1 {
		printHelp()
		return
	}

	queryName := args[0]
	_, ok := conf.Queries[queryName]

	if !ok {
		termRed.Printf("The query %s does not exist.\n", queryName)
	}

	delete(conf.Queries, queryName)
	err := saveConfiguration(conf)

	if err != nil {
		termRed.Printf("delete command: %v\n", err)
	} else {
		termGreen.Printf("Query %s deleted.\n", queryName)
	}
}

func getFromCache(query string) (*CacheFile, error) {
	cache, err := loadCache(query)
	return cache, err
}

func listCommand(args []string, conf *Configuration) {
	if len(args) != 1 {
		printHelp()
		return
	}

	queryName := args[0]
	query, ok := conf.Queries[queryName]
	var result *IssueQueryResult
	var err error

	if !ok {
		termRed.Printf("Query %s is not configured in .issuesrc\n", queryName)
		return
	}

	cache, err := getFromCache(queryName)
	resultFromCache := false

	if err == nil {
		duration := conf.CacheDuration()
		cacheLife := time.Now().Sub(cache.Updated)

		if cacheLife < duration {
			result = &cache.Result
			resultFromCache = true
		}
	}

	if !resultFromCache {
		result, err = queryIssues(strings.Split(query, " "))
	}

	if err != nil {
		termRed.Printf("list command: %v\n", err)
		return
	}

	r := PrintResult{resultFromCache, result}
	r.print()

	if !resultFromCache {
		err = saveCache(queryName, result)
	}
}

func queryCommand(args []string, conf *Configuration) {
	result, err := queryIssues(args)

	if err != nil {
		log.Fatal(err)
	}

	r := PrintResult{false, result}
	r.print()
}

func queriesCommand(args []string, conf *Configuration) {
	termCyan.Printf("Query list:\n")
	for k, v := range conf.Queries {
		fmt.Printf("%s\t%s\n", termBlue.Quote(k), termGreen.Quote(v))
	}
}

func parseArguments(args []string, conf *Configuration) {

	if len(args) < 2 || args[1] == "help" {
		printHelp()
		return
	}

	var commandArg = args[1]
	var arguments = args[2:]
	switch commandArg {
	case "add":
		addCommand(arguments, conf)
	case "delete":
		deleteCommand(arguments, conf)
	case "list":
		listCommand(arguments, conf)
	case "query":
		queryCommand(arguments, conf)
	case "queries":
		queriesCommand(arguments, conf)
	default:
		printHelp()
	}
}

func formatHelpLine(t, d, e string) string {
	return fmt.Sprintf("%s\t%s\n\t%s\n", termBlue.Quote(t), termGreen.Quote(d), e)
}

func printHelp() {
	fmt.Println(termCyan.Quote("Issues - Usage:\n"))
	fmt.Println(formatHelpLine(helpAdd, helpAddDetail, helpAddExample))
	fmt.Println(formatHelpLine(helpDelete, helpDeleteDetail, helpDeleteExample))
	fmt.Println(formatHelpLine(helpList, helpListDetail, helpListExample))
	fmt.Println(formatHelpLine(helpQuery, helpQueryDetail, helpQueryExample))
	fmt.Println(formatHelpLine(helpQueries, helpQueriesDetail, helpQueriesExample))
	fmt.Println(formatHelpLine(helpHelp, helpHelpDetail, helpHelpExample))
}

func main() {

	if !configurationExists() {
		createConfiguration()
	}

	conf, err := loadConfiguration()

	if err != nil {
		panic(termRed.Quote(fmt.Sprintf("%v", err)))
	}

	parseArguments(os.Args, conf)
}
