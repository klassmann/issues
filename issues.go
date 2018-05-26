// Issues helps to search and follow repositories issues list from command line
// You can configure .issuescmdrc to keep saved repositories and predefined searches
// Inspired by https://github.com/adonovan/gopl.io/tree/master/ch4/issues
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	helpAdd        string = "add"
	helpAddDetail  string = "Adds a new query to your local configuration. You need to give a name before the query:"
	helpAddExample string = `$ issues add [query-name] [query parameters]
	$ issues add golang repo:golang/go is:open memory`
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

func printQueryIssues(result *IssueQueryResult) {
	termCyan.Printf("%d issues:\n", result.TotalCount)

	for _, item := range result.Items {

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
			termRed.Printf("add command: %v", err)
		} else {
			termGreen.Printf("Query %s saved.\n", queryName)
		}
	} else {
		printHelp()
	}
}

func listCommand(args []string, conf *Configuration) {
	if len(args) == 1 {
		queryName := args[0]
		query, ok := conf.Queries[queryName]
		if ok {
			result, err := queryIssues(strings.Split(query, " "))

			if err != nil {
				termRed.Printf("list command: %v", err)
			}

			printQueryIssues(result)
		} else {
			termRed.Printf("Query %s is not configured in .issuesrc", queryName)
		}
	} else {
		printHelp()
	}
}

func queryCommand(args []string, conf *Configuration) {
	result, err := queryIssues(args)

	if err != nil {
		log.Fatal(err)
	}

	printQueryIssues(result)
}

func queriesCommand(args []string, conf *Configuration) {
	fmt.Println("Query list:")
	for k, v := range conf.Queries {
		fmt.Printf("\t%s\t%s\n", termBlue.Quote(k), termGreen.Quote(v))
	}
}

func parseArguments(args []string, conf *Configuration) {

	if len(args) < 2 || args[1] == "help" {
		//fmt.Println(commandLineHelp)
		printHelp()
	} else {
		var commandArg = args[1]
		var arguments = args[2:]
		switch commandArg {
		case "add":
			addCommand(arguments, conf)
		case "list":
			listCommand(arguments, conf)
		case "query":
			queryCommand(arguments, conf)
		case "queries":
			queriesCommand(arguments, conf)
		}
	}
}

func formatHelpLine(t, d, e string) string {
	return fmt.Sprintf("%s\t%s\n\t%s\n", termBlue.Quote(t), termGreen.Quote(d), e)
}

func printHelp() {
	fmt.Println(termCyan.Quote("Issues - Usage:\n"))
	fmt.Println(formatHelpLine(helpAdd, helpAddDetail, helpAddExample))
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
