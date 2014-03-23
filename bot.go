package main

import (
	"./configure"
	"./github"

	"github.com/codegangsta/martini"
)

func main() {
	// initializations
	configure.ReadConfiguration()

	m := martini.Classic()
	m.Post("/github_action", github.HandleHook)

	m.Run()
}
