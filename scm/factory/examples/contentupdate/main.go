package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

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

	ctx := context.Background()
	args := os.Args
	if len(args) < 7 {
		fmt.Printf("ar8uments: owner repository remotePath localPath branch commitMsg shaOfOldContent\n")
		return
	}

	owner := args[1]
	repo := args[2]
	remotePath := args[3]
	localPath := args[4]
	branch := args[5]
	commitMsg := args[6]
	sha := args[7]

	data, err := ioutil.ReadFile(localPath) // #nosec
	if err != nil {
		fmt.Printf("unable to load file from localPath : %v", err)
		return
	}

	fullRepo := scm.Join(owner, repo)
	cp := &scm.ContentParams{
		Ref:     "",
		Branch:  branch,
		Message: commitMsg,
		Data:    data,
		Sha:     sha,
	}

	fmt.Printf("updating content for repository %s/%s and remotePath: %s with branch: %s\n", owner, repo, remotePath, branch)
	_, err = client.Contents.Update(ctx, fullRepo, remotePath, cp)
	if err != nil {
		helpers.Fail(err)
		return
	}

}

func createListOptions() scm.ListOptions {
	return scm.ListOptions{
		Size: 1000,
	}
}
