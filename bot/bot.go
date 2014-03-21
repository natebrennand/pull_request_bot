package main

import (
	"./github"

	"fmt"
	"github.com/codegangsta/martini"
)

func main() {
	config := Configure()
	fmt.Println(config)

	m := martini.Classic()
	m.Post("/github_action", github.HandleHook)

	m.Run()
}
