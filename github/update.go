package github

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PullRequestResponse struct {
	Sha     string
	Merged  bool
	Message string
}

type MergeRequestStruct struct {
	CommitMessage string `json:"commit_message"`
}

func GithubAPICall(token, uri, method string, data []byte) ([]byte, error) {
	uri = "https://api.github.com" + uri
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func MergePullRequest(token, owner, repo string, requestNumber int) error {
	uri := ("/repos/" + repo + "/pulls/" + strconv.Itoa(requestNumber) + "/merge")
	bodyBytes, err := json.Marshal(MergeRequestStruct{"merge bot!"})
	if err != nil {
		return err
	}
	_, err = GithubAPICall(token, uri, "PUT", bodyBytes)
	return err
}
