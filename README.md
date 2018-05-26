# Issues Command Line

This small utility helps query and follow repositories issues.

## Requirements

- Golang


## How to install

    $ go install github.com/klassmann/issues

## Usage

    add	Adds a new query to your local configuration. You need to give a name to the query:
            $ issues add [query-name] [query parameters]
            $ issues add golang repo:golang/go is:open memory

    delete	Deletes a query from your local configuration.
            $ issues delete [query-name]
            $ issues delete golang

    list	Query issues from a [query-name] saved in local configuration:
            $ issues list [query-name]

    query	Queries and shows a new search on github:
            $ issues query [query-parameters]
            $ issues query repo:golang/go is:open memory

    queries	Show the list of queries in your configuration.
            $ issues queries

    help	Shows this help message.


## Configuration
Create a `.issuesrc` in your home directory. The file will be created automatically on first use.

    {
        "queries": 
        {
            "sarama": "repo:shopify/sarama is:open",
            "golang": "repo:golang/go is:open",
            "django": "repo:django/django is:open"
        }
    }

Note: Each query should have a name and a GitHub issue search string, see more [here](https://developer.github.com/v3/search/#search-issues).


## Notes

- Inspired by [The Go Programming Language - Issue example](https://github.com/adonovan/gopl.io/tree/master/ch4/issues)


## License
[Apache 2.0](LICENSE)
