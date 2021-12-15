package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/pflag"
)

const (
	githubTrafficClonesURLTemplate = "https://api.github.com/repos/%s/%s/traffic/clones"
	trafficClonesLabel             = "clones"
)

type trafficClonesResp struct {
	Count   int           `json:"count"`
	Uniques int           `json:"uniques"`
	Clones  []interface{} `json:"clones"`
}

type req struct {
	GitUser string `form:"git_user" binding:"required"`
	GitRepo string `form:"git_repo" binding:"required"`
}

// reference: https://shields.io/endpoint
type response struct {
	SchemaVersion int    `json:"schemaVersion"`
	Label         string `json:"label"`
	Message       string `json:"message"`
	Color         string `json:"color"`
}

type errResponse struct {
	Message string `json:"message"`
}

func main() {
	var githubToken string
	pflag.StringVarP(&githubToken, "token", "t", "", "Your github personal access token(https://github.com/settings/tokens)")
	pflag.Parse()

	if githubToken == "" {
		fmt.Fprint(os.Stderr, "must provide your github token by flag --token/-t")
		os.Exit(1)
	}

	r := gin.Default()

	r.GET("/v1/repo-traffic-clones", func(c *gin.Context) {
		var r req
		if err := c.ShouldBindQuery(&r); err != nil {
			c.JSON(400, errResponse{
				Message: err.Error(),
			})
			return
		}

		var trafficClonesResp trafficClonesResp

		request := gorequest.New().Get(fmt.Sprintf(githubTrafficClonesURLTemplate, r.GitUser, r.GitRepo))
		request.Set("Authorization", "token "+githubToken)

		resp, _, errs := request.EndStruct(&trafficClonesResp)
		if len(errs) != 0 {
			c.JSON(500, errResponse{
				Message: fmt.Sprintf("github api resp: %v, errors: %v", resp, errs),
			})
			return
		}

		c.JSON(200, response{
			SchemaVersion: 1,
			Label:         trafficClonesLabel,
			Message:       strconv.Itoa(trafficClonesResp.Count),
			Color:         "orange",
		})
	})

	_ = r.Run(":9000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
