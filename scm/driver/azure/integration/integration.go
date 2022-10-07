package integration

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/jenkins-x/go-scm/scm"
)

var (
	client *scm.Client
	token  = GetToken()

	// The name of the Azure DevOps organization.
	organization = "tphoney"

	// Project ID or project name
	project = "test_project"

	// The name or ID of the repository.
	// This repo should have a single file in, with a branch called pr_branch that can be used to create a PR from. e.g.
	// ❯ tree
	// .
	// └── README.md
	//
	// 0 directories, 1 file
	repositoryID = "fde2d21f-13b9-4864-a995-83329045289a"
)

func GetCurrentCommitOfBranch(client *scm.Client, branch string) (string, error) {
	commit, _, err := client.Contents.List(context.Background(), repositoryID, "", "main")
	if err != nil {
		return "", err
	}
	return commit[0].Sha, nil
}

func GetToken() string {
	rawToken := os.Getenv("AZURE_TOKEN")
	if rawToken != "" {
		return base64.StdEncoding.EncodeToString([]byte(":" + rawToken))
	}
	return ""
}
