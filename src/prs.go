package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type PullRequestCommit struct {
	SHA         *string         `json:"sha,omitempty"`
	NodeID      *string         `json:"node_id,omitempty"`
	Commit      *github.Commit  `json:"commit,omitempty"`
	Author      *github.User    `json:"author,omitempty"`
	Committer   *github.User    `json:"committer,omitempty"`
	Parents     []github.Commit `json:"parents,omitempty"`
	HTMLURL     *string         `json:"html_url,omitempty"`
	URL         *string         `json:"url,omitempty"`
	CommentsURL *string         `json:"comments_url,omitempty"`
}

func ListCommits(pullRequest *github.PullRequest) ([]*PullRequestCommit, error) {
	var commits []*PullRequestCommit
	r, err := httpClient.Get(*pullRequest.URL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	decodeErr := json.NewDecoder(r.Body).Decode(&commits)

	return commits, decodeErr
}

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := github.PullRequestListOptions{State: "closed"}

	prs, _, err := client.PullRequests.List(ctx, "apache", "incubator-superset", &opts)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("pull requests:")
		fmt.Println()
		for index, pr := range prs {
			fmt.Printf("%v.\t[%v] %v\n", index, *pr.ID, *pr.Title)
			fmt.Printf("\t%v\n", *pr.URL)

			commits, err := ListCommits(pr)
			if err != nil {
				fmt.Printf("\t\terror while fetching commits: %v\n", err)
			} else {
				if pr.Commits != nil {
					fmt.Printf("\t\tcommits(%v):\n", *pr.Commits)
					for _, commit := range commits {
						fmt.Printf("\t\t[%v] %v\n", *commit.SHA, *commit.URL)
					}
				}
			}
			fmt.Printf("\tby %v\n", *pr.User.Login)
			fmt.Println()
		}
	}
}
