package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

const (
	typeClones = "clones"
)

type req struct {
	GitUser string `form:"git_user" binding:"required"`
	GitRepo string `form:"git_repo" binding:"required"`
	// "count" or "uniques"
	Type  string `form:"type" binding:"required"`
	Label string `form:"label" binding:"required"`
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
	pflag.StringVarP(
		&githubToken,
		"token",
		"t",
		"",
		"Your github personal access token(https://github.com/settings/tokens)",
	)
	pflag.Parse()

	if githubToken == "" {
		fmt.Fprint(os.Stderr, "must provide your github token by flag --token/-t")
		os.Exit(1)
	}

	router := gin.Default()

	router.GET("/v1/github/traffic/clones/total", func(c *gin.Context) {
		var r req
		if err := c.ShouldBindQuery(&r); err != nil {
			c.JSON(400, errResponse{
				Message: err.Error(),
			})
			return
		}

		githubClones, err := getGithubTrafficClones(r.GitUser, r.GitRepo, githubToken)
		if err != nil {
			c.JSON(500, errResponse{
				Message: err.Error(),
			})

			return
		}

		updateGithubTrafficClones(githubClones.Clones, r.GitUser, r.GitRepo)

		total := getClonesTotal(r.GitUser, r.GitRepo)

		var typeTotal int
		if r.Type == "count" {
			typeTotal = total.Count
		} else if r.Type == "uniques" {
			typeTotal = total.Uniques
		} else {
			c.JSON(400, errResponse{
				Message: "type must be 'count' or 'uniques'",
			})

			return
		}

		c.JSON(200, response{
			SchemaVersion: 1,
			Label:         r.Label,
			Message:       strconv.Itoa(typeTotal),
			Color:         "orange",
		})
	})

	_ = router.Run(":9000")
}
