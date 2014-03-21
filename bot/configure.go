package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Repo struct {
	Name            string
	Approvers       []string
	ApprovalsNeeded int
}

// open the file and read the data into the Repo struct
func FromJson(rl *[]Repo, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Config file, '%s', not found.\n", filename)
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(rl)
	if err != nil {
		fmt.Printf("Problems reading your config file\n%s\n", err.Error())
		return err
	}

	return nil
}

// gathers variables from the environment
func GetEnvVariables() map[string]string {
	envVars := make(map[string]string)

	for _, varStr := range os.Environ() {
		split := strings.Split(varStr, "=")
		envVars[split[0]] = split[1]
	}
	return envVars
}

// returns an array of Repo structs
func Configure() []Repo {
	configFile := GetEnvVariables()["config"]
	if len(configFile) == 0 {
		configFile = "config.json"
	}

	var repos []Repo
	err := FromJson(&repos, configFile)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	return repos
}
