package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Action struct {
	Repository struct {
		Name      string
		Full_name string
		Owner     User
	}
	Sender  User
	Pull_request PullRequest
	Comment Comment
	Action	string
}

type PullRequest struct {
	Title     string
	Body      string
	Html_url  string
	Merged    bool
	Mergeable bool
	Comments  int
}

type Comment struct {
	Body     string
	Html_url string
	User     User
}

type User struct {
	Login string
	Id    int
}

// read in request data to a Action struct
func ParseData(req *http.Request) (Action, error) {
	// parse bytes from req object
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err.Error())
		return Action{}, err
	}

	// build struct from byte array
	var body Action
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return Action{}, err
	}

	// return new Action struct
	return body, nil
}

func HandleHook(req *http.Request) (int, string) {
	body, err := ParseData(req)
	if err != nil {
		return 500, "SERVER ERROR: " + err.Error()
	}
	fmt.Printf("json: %#v\n", body)
	fmt.Printf("json: %#v\n", body.Repository)

	return 200, "RECEIVED!"
}
