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

func printQueryIssues(result *IssueQueryResult) {
	fmt.Printf("%d issues:\n", result.TotalCount)

	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.80s\n",
			item.Number, item.User.Login, item.Title)
	}
}

func addCommand(args []string, conf *Configuration) {

}

func listCommand(args []string, conf *Configuration) {
	if len(args) == 1 {
		queryName := args[0]
		query, ok := conf.Queries[queryName]
		if ok {
			result, err := queryIssues(strings.Split(query, " "))

			if err != nil {
				log.Fatal(err)
			}

			printQueryIssues(result)
		} else {
			fmt.Printf("Query %s is not configured in .issuesrc", queryName)
		}
	} else {
		fmt.Println(commandLineHelp)
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

}

func parseArguments(args []string, conf *Configuration) {

	if len(args) < 2 || args[1] == "help" {
		fmt.Println(commandLineHelp)
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

func main() {
	if !configurationExists() {
		createConfiguration()
	}

	conf, err := loadConfiguration()

	if err != nil {
		panic(err)
	}

	parseArguments(os.Args, conf)
}
