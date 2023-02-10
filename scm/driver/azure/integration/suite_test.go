package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

var (
	// personal access token that can reach the organization
	token        = os.Getenv("AZURE_TOKEN")
	organization = os.Getenv("AZURE_ORG")
	// this project should be safe to create/delete repositories in
	project   = os.Getenv("AZURE_PROJECT")
	projectFQ = fmt.Sprintf("%s/%s", organization, project)
	client    = makeClient()
)

type TestCase struct {
	Name string
	Test func(t *testing.T)
}

func canRun(t *testing.T) {
	if (token != "") && (project != "") && (organization != "") {
		return
	}
	t.Skip("Acceptance tests not configured (need AZURE_TOKEN, AZURE_ORG, AZURE_PROJECT)")
}

func repoFQ(repositoryName string) string {
	return fmt.Sprintf("%s/%s", projectFQ, repositoryName)
}

func makeClient() *scm.Client {
	cl, err := factory.NewClient("azure", "https://dev.azure.com", token)
	if err != nil {
		panic("could not create azure client")
	}
	return cl
}

func makeCleanRepo(name string) (*scm.Repository, error) {
	// clean old one, drop the error
	_, _ = client.Repositories.Delete(context.Background(), repoFQ(name))

	ri := scm.RepositoryInput{
		Namespace: projectFQ,
		Name:      name,
	}

	testRepo, _, err := client.Repositories.Create(context.Background(), &ri)
	return testRepo, err
}

func initializeRepo(cloneURL string) (string, error) {
	// Using go-git because it avoids a chicken and egg problem with Contents needing
	// to work to set up the test repo.

	fs := memfs.New()
	r, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		return "", fmt.Errorf("failed to create in memory repository: %w", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return "", fmt.Errorf("could not get working tree: %w", err)
	}

	readme, err := fs.Create("README.md")
	if err != nil {
		return "", fmt.Errorf("could not open README.md for writing: %w", err)
	}

	_, err = readme.Write([]byte("# Content Test"))
	if err != nil {
		return "", fmt.Errorf("could not write to README.md: %w", err)
	}

	err = readme.Close()
	if err != nil {
		return "", fmt.Errorf("could not close README.md: %w", err)
	}

	_, err = w.Add("README.md")
	if err != nil {
		return "", fmt.Errorf("could not add README.md to the worktree: %w", err)
	}

	h, err := w.Commit("initialised with README", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Author McAuthorson",
			Email: "amca@example.com",
			When:  time.Now(),
		},
		Parents: nil,
	})
	if err != nil {
		return "", fmt.Errorf("could not create initial commit: %w", err)
	}
	readmeCommitSha := h.String()

	auth := &gitHttp.BasicAuth{
		Username: "",
		Password: token,
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{cloneURL},
	})

	if err != nil {
		return readmeCommitSha, fmt.Errorf("could not create the remote: %w", err)
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		RemoteURL:  cloneURL,
		RefSpecs:   []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:refs/heads/main", readmeCommitSha))},
		Auth:       auth,
	})
	if err != nil {
		return readmeCommitSha, fmt.Errorf("could not push repo: %w", err)
	}

	return readmeCommitSha, nil
}
