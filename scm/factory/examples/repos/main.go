package main

import (
	"context"
	"fmt"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/factory/examples/helpers"
)

func main() {
	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		helpers.Fail(err)
		return
	}

	fmt.Printf("listing repostories\n")

	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, createListOptions())
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("Found %d repostories\n", len(repos))

	for _, r := range repos {
		fmt.Printf("  repo: %#v\n", r)
	}
}

func createListOptions() scm.ListOptions {
	return scm.ListOptions{
		Size: 1000,
	}
}
