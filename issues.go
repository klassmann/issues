// Issues helps to search and follow repositories issues list from command line
// You can configure .issuescmdrc to keep saved repositories and predefined searches
// Inspired by https://github.com/adonovan/gopl.io/tree/master/ch4/issues
package main

import (
	"fmt"
	"log"
	"os"
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
		fmt.Printf("#%-5d %9.9s %.80s\n",
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
	if !configurationExists() {
		createConfiguration()
	}
	parseArguments(os.Args)
}
