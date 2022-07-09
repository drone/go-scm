package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/jenkins-x/go-scm/scm/factory"
	"github.com/jenkins-x/go-scm/scm/factory/examples/helpers"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("usage: org/repo prNumber [label]")
		return
	}
	repo := args[1]

	client, err := factory.NewClientFromEnvironment()
	if err != nil {
		helpers.Fail(err)
		return
	}
	ctx := context.Background()

	prText := args[2]
	number, err := strconv.Atoi(prText)
	if err != nil {
		helpers.Fail(errors.Wrapf(err, "failed to parse PR number: %s", prText))
		return
	}

	fmt.Printf("Getting PR\n")

	pr, _, err := client.PullRequests.Find(ctx, repo, number)
	if err != nil {
		helpers.Fail(err)
		return
	}
	fmt.Printf("Found PullRequest:\n")
	data, err := yaml.Marshal(pr)
	if err != nil {
		helpers.Fail(errors.Wrap(err, "failed to marshal PR as YAML"))
		return
	}
	fmt.Printf("%s:\n", string(data))

	labels, _, err := client.PullRequests.ListLabels(ctx, repo, number, &scm.ListOptions{})
	if err != nil {
		helpers.Fail(errors.Wrap(err, "failed to list labels on PR"))
		return
	}

	for _, l := range labels {
		fmt.Printf("label %s\n", l.Name)
	}

	fmt.Printf("done")
}
