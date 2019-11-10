# Pinboard

[![GoDoc](https://godoc.org/github.com/drags/pinboard?status.svg)](https://godoc.org/github.com/drags/pinboard)

Golang client library for the Pinboard

Get Started:

 * Install with `go get github.com/drags/pinboard`
 * Browse [![GoDoc](https://godoc.org/github.com/drags/pinboard?status.svg)](https://godoc.org/github.com/drags/pinboard)
 * Browse [Pinboard API](https://pinboard.com/api) documentation


## Example

```go
package main

import (
	"fmt"
	"github.com/drags/pinboard"
)

func main() {
	p := pinboard.Pinboard{User: "myuser", Token: "6D3F1921A2C33D82EA1"}

	// Get the 25 most recent posts
	prf := pinboard.PostsRecentFilter{}
	prf.Count = 25
	posts, err := p.PostsRecent(prf)
	if err != nil {
		fmt.Println("Got error from PostsRecent", err)
	}
	for _, v := range posts {
		fmt.Println("Post URL", v.Url)
		fmt.Println("Post description", v.Description)
		fmt.Println("Post tags", v.Tags)
	}

	// Add a new post
	pp := pinboard.Post{
		Url:         "https://example.com",
		Description: "This field should actually be called title, but backward compat is a thing",
		Tags:        []string{"tag1", ".private-tag"},
		Extended:    "The actual body of the post",
	}

	err = p.PostsAdd(pp, false, false)
	if err != nil {
		fmt.Println("Got error from PostsAdd", err)
	}
}
```

## Installation

To install Pinboard:
```
go get github.com/drags/pinboard
```

## Staying up to date

To update Pinboard to the latest version, use `go get -u github.com/drags/pinboard`.

## Bug reports

Report bugs using the [drags/pinboard#Issues](https://github.com/drags/pinboard/issues) tracker.
