package github

import (
	"../configure"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"container/list"
)

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
	Sender  User
	Comment Comment `json:"comment"`
	Action  string  `json:"action"`
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

// checks if a comment is approving the pull request
func (c Comment) RequestApproved(approvers []string, approvals *list.List) bool {
	for _, approver := range approvers {
		// check for users already used
		for a := approvals.Front(); a != nil; a = a.Next() {
			if a.Value == approver {
				// nil the approver so it's not matched up
				approver = ""
			}
		}

		// check if user is approved
		if approver == c.User.Login {
			for _, keyword := range configure.GlobalConfig.MergePhrases {
				if strings.Contains(c.Body, keyword) {
					approvals.PushBack(c.User.Login)
					return true
				}
			}
		}
	}
	return false
}

// given a set of comments, determine if the request is ready to be merged
func CheckIssue(owner, repo string, comments []Comment) bool {
	config := configure.GlobalConfig.Repos
	approvalsNeeded := 0
	approvers := []string{}
	madeApproval := list.New()

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
		if comment.RequestApproved(approvers, madeApproval) {
			approvalsNeeded -= 1
		}
	}

	// if not enough approvals have been made yet
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

// Checks the comments on an issue
func CheckIssueComments(owner, repo string, issueNumber int) error {
	uri := "/repos/" + repo + "/issues/" + strconv.Itoa(issueNumber) + "/comments"
	respBody, err := GithubAPICall(uri, "GET", nil)
	if err != nil {
		return err
	}
	var comments []Comment
	json.Unmarshal(respBody, &comments)

	if CheckIssue(owner, repo, comments) {
		return MergePullRequest(owner, repo, issueNumber)
	}
	return nil
}

// recieves the webhook from github
func HandleHook(req *http.Request) (int, string) {
	body, err := ParseData(req)
	if err != nil {
		return 500, "SERVER ERROR: " + err.Error()
	}

	if body.Action == "created" {
		err = CheckIssueComments(body.Repository.Owner.Login, body.Repository.FullName, body.Issue.Number)
		if err != nil {
			return 400, err.Error()
		}
	}

	return 200, "RECEIVED!"
}
