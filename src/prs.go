package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func main() {
	ctx := context.Background()
	client := github.NewClient(nil)

	opts := github.PullRequestListOptions{State: "closed"}

	prs, _, err := client.PullRequests.List(ctx, "apache", "incubator-superset", &opts)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("Pull Requests:")
		fmt.Println()
		for index, pr := range prs {
			fmt.Printf("%v.\t[%v] %v\n", index, *pr.ID, *pr.Title)
			fmt.Printf("\t%v\n", *pr.URL)
			if pr.Commits != nil {
				fmt.Printf("\t[%v] %v\n", *pr.Commits, *pr.CommitsURL)
			}
			fmt.Printf("\tby %v\n", *pr.User.Login)
			fmt.Println()
		}
	}
}
