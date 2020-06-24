package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/factory/examples/helpers"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: repo branch")
		os.Exit(1)
		return
	}
	repo := args[1]
	ref := "master"
	if len(args) > 2 {
		ref = args[2]
	}

	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		helpers.Fail(err)
		return
	}

	fmt.Printf("finding in repo: %s ref: %s\n", repo, ref)

	ctx := context.Background()
	commits, _, err := client.Git.ListCommits(ctx, repo, scm.CommitListOptions{
		Ref:  ref,
		Page: 1,
		Size: 100,
	})
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("found %d commits\n", len(commits))
	for _, commit := range commits {
		fmt.Printf("commit %s by %s message: %s\n", commit.Sha, commit.Committer.Name, commit.Message)
	}
}
