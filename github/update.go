package github

import (
	"../configure"

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

func GithubAPICall(uri, method string, data []byte) ([]byte, error) {
	uri = "https://api.github.com" + uri
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+configure.GlobalConfig.GithubToken)

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

func MergePullRequest(owner, repo string, requestNumber int) error {
	uri := ("/repos/" + repo + "/pulls/" + strconv.Itoa(requestNumber) + "/merge")
	bodyBytes, err := json.Marshal(MergeRequestStruct{"merge bot merging!"})
	if err != nil {
		return err
	}
	_, err = GithubAPICall(uri, "PUT", bodyBytes)
	return err
}
