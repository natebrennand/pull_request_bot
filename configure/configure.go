package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var GlobalConfig UtilConfig

type UtilConfig struct {
	GithubToken string
	Repos []Repo
}

type Repo struct {
	Name            string
	Approvers       []string
	ApprovalsNeeded int
}

// open the file and read the data into the Repo struct
func FromJson(c *UtilConfig, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Config file, '%s', not found.\n", filename)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
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

// gathers configuration and sets the module variable "GlobalConfig"
func Configure() {
	configFile := GetEnvVariables()["config"]
	if len(configFile) == 0 {
		configFile = "config.json"
	}

	var c UtilConfig
	err := FromJson(&c, configFile)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// set module var
	GlobalConfig = c
}
