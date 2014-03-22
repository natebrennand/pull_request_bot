package github

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func MergePullRequest(token, repo, owner string, requestNumber int) error {
	uri := ("https://api.github.com/repos/" + owner + "/" + repo + "/pulls/" + strconv.Itoa(requestNumber) + "/merge")
	fmt.Println(uri)

	bodyBytes, err := json.Marshal(MergeRequestStruct{"merge bot!"})
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest("PUT", uri, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))
	return nil
}
