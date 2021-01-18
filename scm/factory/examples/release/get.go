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
	if len(args) < 3 {
		fmt.Println("usage: repo tag")
		os.Exit(1)
		return
	}
	repo := args[1]
	tag := args[2]
	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		helpers.Fail(err)
		return
	}

	fmt.Printf("finding release for repo: %s tag: %s\n", repo, tag)

	ctx := context.Background()
	release, _, err := client.Releases.FindByTag(ctx, repo, tag)
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("found %v\n", release)
}
