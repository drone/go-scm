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

	fmt.Printf("finding in repo: %s ref: %s\n", repo, ref)

	ctx := context.Background()
	changes, _, err := client.Git.ListChanges(ctx, repo, ref, scm.ListOptions{
		Page: 1,
		Size: 100,
	})
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("found %d changes\n", len(changes))
	for _, change := range changes {
		action := "changed"
		if change.Added {
			action = "added"
		} else if change.Deleted {
			action = "deleted"
		} else if change.Renamed {
			action = "renamed"
		}
		fmt.Printf("%s on path %s at: %s\n", action, change.Path, change.Sha)
	}
}
