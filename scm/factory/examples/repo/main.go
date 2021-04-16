package main

import (
	"context"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"os"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/factory/examples/helpers"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Printf("usage: repoName")
		return
	}

	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		helpers.Fail(err)
		return
	}
	ctx := context.Background()

	name := args[1]
	fmt.Printf("finding repository %s\n", name)

	repo, _, err := client.Repositories.Find(ctx, name)
	if scm.IsScmNotFound(err) {
		fmt.Printf("not found\n")
		return
	}
	if err != nil {
		fmt.Printf("failed: %s\n", err.Error())
		return
	}
	fmt.Printf("found %#v\n", repo)
}
