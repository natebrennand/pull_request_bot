package main

import (
	"./github"
	"./configure"

	"github.com/codegangsta/martini"
)

func main() {
	config := configure.Configure()

	m := martini.Classic()
	m.Post("/github_action", github.HandleAction)

	m.Run()
}
