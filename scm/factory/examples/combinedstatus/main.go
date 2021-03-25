package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/factory/examples/helpers"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: repo ref")
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

	fmt.Printf("finding combined status in repo: %s ref: %s\n", repo, ref)

	ctx := context.Background()
	results, _, err := client.Repositories.FindCombinedStatus(ctx, repo, ref)
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("state: %v sha: %s\n", results.State, results.Sha)

	for _, s := range results.Statuses {
		fmt.Printf("target %s state %v label %s\n", s.Target, s.State, s.Label)
	}
}
