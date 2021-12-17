package main

import (
	"fmt"

	"github.com/parnurzeal/gorequest"
)

const githubTrafficClonesURLTemplate = "https://api.github.com/repos/%s/%s/traffic/clones"

type trafficClonesResp struct {
	Count   int          `json:"count"`
	Uniques int          `json:"uniques"`
	Clones  []clonesItem `json:"clones"`
}

type clonesItem struct {
	Timestamp string `json:"timestamp"`
	Count     int    `json:"count"`
	Uniques   int    `json:"uniques"`
}

func getGithubTrafficClones(user, repoName, token string) (*trafficClonesResp, error) {
	var trafficClonesResp trafficClonesResp

	request := gorequest.New().Get(fmt.Sprintf(githubTrafficClonesURLTemplate, user, repoName))
	request.Set("Authorization", "token "+token)

	resp, _, errs := request.EndStruct(&trafficClonesResp)
	if len(errs) != 0 {
		return nil, fmt.Errorf("github api resp: %v, errors: %v", resp, errs)
	}

	return &trafficClonesResp, nil
}
