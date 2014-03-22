package github

import (
	"../configure"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var KEYWORDS = [2]string{"lgtm", "merge it"}

type Action struct {
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    User
	}
	Issue struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	Sender      User
	PullRequest PullRequest `json:"pull_request"`
	Comment     Comment     `json:"comment"`
	Action      string      `json:"action"`
}

type PullRequest struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	HtmlUrl   string `json:"html_url"`
	Merged    bool   `json:"merged"`
	Mergeable bool   `json:"mergeable"`
	Comments  int    `json:"comments"`
}

type Comment struct {
	Body    string `json:"body"`
	HtmlUrl string `json:"html_url"`
	User    User
}

type User struct {
	Login string `json:"login"`
	Id    int    `json:"id"`
}

func (c Comment) RequestApproved(approvers []string) bool {
	for _, approver := range approvers {
		if approver == c.User.Login {
			for _, keyword := range KEYWORDS {
				if strings.Contains(c.Body, keyword) {
					return true
				}
			}
		}
	}
	return false
}

func CheckIssue(owner, repo string, comments []Comment, config []configure.Repo) bool {
	approvalsNeeded := 0
	approvers := []string{}

	for _, relevantRepo := range config {
		if relevantRepo.Name == repo {
			approvalsNeeded = relevantRepo.ApprovalsNeeded
			approvers = relevantRepo.Approvers
		}
	}
	if approvalsNeeded == 0 { // repo not in config
		return false
	}

	for _, comment := range comments {
		if comment.RequestApproved(approvers) {
			approvalsNeeded -= 1
		}
	}

	if approvalsNeeded > 0 {
		return false
	}
	return true
}

// read in request data to a Action struct
func ParseData(req *http.Request) (Action, error) {
	// parse bytes from req object
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
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

func CheckIssueComments(token, owner, repo string, issueNumber int, config []configure.Repo) error {
	uri := "/repos/" + repo + "/issues/" + strconv.Itoa(issueNumber) + "/comments"
	respBody, err := GithubAPICall(token, uri, "GET", nil)
	if err != nil {
		return err
	}
	var comments []Comment
	json.Unmarshal(respBody, &comments)

	if CheckIssue(owner, repo, comments, config) {
		return MergePullRequest(token, owner, repo, issueNumber)
	}
	return nil
}

func HandleHook(req *http.Request) (int, string) {
	body, err := ParseData(req)
	if err != nil {
		return 500, "SERVER ERROR: " + err.Error()
	}

	config := configure.GlobalConfig
	if body.Action == "created" {
		err = CheckIssueComments(configure.GlobalEnv["github_token"], body.Repository.Owner.Login, body.Repository.FullName, body.Issue.Number, config)
		if err != nil {
			return 400, err.Error()
		}
	}

	return 200, "RECEIVED!"
}
