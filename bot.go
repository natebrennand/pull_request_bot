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
	m.Post(configure.GlobalConfig.WebhookURI, github.HandleHook)

	m.Run()
}
