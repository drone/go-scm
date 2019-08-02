package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

func main() {
	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		fail(err)
		return
	}

	fmt.Printf("listing repostories\n")

	ctx := context.Background()
	repos, _, err := client.Repositories.List(ctx, createListOptions())
	if err != nil {
		fail(err)
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

func fail(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
	os.Exit(1)
}
