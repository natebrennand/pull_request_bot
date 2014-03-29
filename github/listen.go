package github

import (
	"../configure"

	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Action struct {
	Repository Repo `json:"repository"`
	Issue      struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		User   User   `json:"user"`
	}
	Sender  User
	Comment Comment `json:"comment"`
	Action  string  `json:"action"`
}

type Repo struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User
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
func (c Comment) RequestApproved(approvers []string, used *list.List) bool {
	for _, approver := range approvers {
		// check for users already used
		for a := used.Front(); a != nil; a = a.Next() {
			if a.Value == approver {
				// nil the approver so it's not matched up
				approver = ""
			}
		}

		// check if user is approved
		if approver == c.User.Login {
			for _, keyword := range configure.GlobalConfig.MergePhrases {
				if strings.Contains(c.Body, keyword) {
					used.PushBack(c.User.Login)
					return true
				}
			}
		}
	}
	return false
}

// given a set of comments, determine if the request is ready to be merged
func (r *Repo) CheckIssue(comments []Comment, sender string) bool {
	config := configure.GlobalConfig.Repos
	approvalsNeeded := 0
	approvers := []string{}
	used := list.New()
	used.PushBack(sender) // block issue creator from approving

	for _, relevantRepo := range config {
		if relevantRepo.Name == r.FullName {
			approvalsNeeded = relevantRepo.ApprovalsNeeded
			approvers = relevantRepo.Approvers
		}
	}
	if approvalsNeeded == 0 { // repo not in config
		return false
	}

	for _, comment := range comments {
		if comment.RequestApproved(approvers, used) {
			approvalsNeeded -= 1
		}
	}

	// if not enough approvals have been made yet
	if approvalsNeeded > 0 {
		return false
	}
	return true
}

// Checks the comments on an issue
func (r *Repo) CheckIssueComments(issueNumber int, sender string) error {
	uri := "/repos/" + r.Owner.Login + "/issues/" + strconv.Itoa(issueNumber) + "/comments"
	respBody, err := GithubAPICall(uri, "GET", nil)
	if err != nil {
		return err
	}
	var comments []Comment
	json.Unmarshal(respBody, &comments)

	if r.CheckIssue(comments, sender) {
		fmt.Println("We are trying to merge PR the corresponding PR\n")
		return MergePullRequest(r.Owner.Login, r.FullName, issueNumber)
	}
	return nil
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

	return body, nil // return new Action struct
}

// recieves the webhook from github
func HandleHook(req *http.Request) (int, string) {
	body, err := ParseData(req)
	if err != nil {
		return 500, "SERVER ERROR: " + err.Error()
	}

	if body.Action == "created" {
		fmt.Printf("%s made a comment on issue #%d\n", body.Sender.Login, body.Issue.Number)
		err = body.Repository.CheckIssueComments(body.Issue.Number, body.Sender.Login)
		if err != nil {
			return 400, err.Error()
		}
	}

	return 200, "RECEIVED!"
}
