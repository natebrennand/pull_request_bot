package main

import (
	"./configure"
	"./github"

	"github.com/codegangsta/martini"
)

func main() {
	// initializations
	_ = configure.Configure()
	_ = configure.GetEnvVariables()

	m := martini.Classic()
	m.Post("/github_action", github.HandleHook)

	m.Run()
}
