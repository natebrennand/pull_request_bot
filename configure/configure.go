package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var GlobalConfig UtilConfig

type UtilConfig struct {
	GithubToken  string
	WebhookURI   string
	MergePhrases []string
	Repos        []Repo
}

type Repo struct {
	Name            string
	Approvers       []string
	ApprovalsNeeded int
}

// open the file and read the data into the Repo struct
func FromJson(c *UtilConfig, filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		fmt.Printf("Config file, '%s', not found.\n", filename)
		return err
	}

	err = json.NewDecoder(file).Decode(c)
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
		splitEnvStatement := strings.Split(varStr, "=")
		envVars[splitEnvStatement[0]] = splitEnvStatement[1]
	}
	return envVars
}

// gathers configuration and sets the module variable, GlobalConfig
func ReadConfiguration() {
	configFilename := GetEnvVariables()["config"]
	if len(configFilename) == 0 {
		configFilename = "config.json"
	}

	err := FromJson(&GlobalConfig, configFilename)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
