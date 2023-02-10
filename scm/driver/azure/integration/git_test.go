package integration_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

func TestGit(t *testing.T) {
	canRun(t)

	const testRepoName = "git-test"
	var testRepo *scm.Repository
	var readmeCommitSha string
	var branchRef *scm.Reference
	setup := func(t *testing.T) {
		var err error
		testRepo, err = makeCleanRepo(testRepoName)
		if err != nil {
			t.Fatalf("could not create repo: %+v", err)
		}

		readmeCommitSha, err = initializeRepo(testRepo.Clone)
		if err != nil {
			t.Fatalf("could not initialize repo: %+v", err)
		}

	}

	tests := []TestCase{
		{
			Name: "create branches",
			Test: func(t *testing.T) {
				var err error
				var res *scm.Response
				branchRef, res, err = client.Git.CreateRef(context.Background(), repoFQ(testRepoName), "branch1", readmeCommitSha)
				if err != nil {
					t.Errorf("could not create ref 'branch1': %v", err)
				}

				// Azure strangely gives a 200
				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}
			},
		},
		{
			Name: "list branches",
			Test: func(t *testing.T) {
				branches, res, err := client.Git.ListBranches(context.Background(), repoFQ(testRepoName), &scm.ListOptions{})
				if err != nil {
					t.Errorf("could not list branches: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if len(branches) != 2 {
					t.Errorf("expected 2 branches (main and branch 1), got %d", len(branches))
				}

				var foundBranch *scm.Reference
				for _, ref := range branches {
					if ref.Name == "branch1" {
						foundBranch = ref
						break
					}
				}

				if foundBranch == nil {
					t.Errorf("could not find 'branch1'")
					return
				}

				if branchRef.Sha != foundBranch.Sha {
					t.Errorf("expected 'branch1' to have sha `%s`, got `%s`", branchRef.Sha, foundBranch.Sha)
				}
			},
		},
		{
			Name: "list commits",
			Test: func(t *testing.T) {
				// add a commit
				content := scm.ContentParams{
					Message: "create a main.go app",
					Data:    []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello world!\")\n}"),
					Branch:  "refs/heads/main",
					Ref:     readmeCommitSha,
				}
				res, err := client.Contents.Create(context.Background(), repoFQ(testRepoName), "/main.go", &content)
				if err != nil {
					t.Errorf("could not create main.go file and commit: %v", err)
				}

				if res.Status != http.StatusCreated {
					t.Errorf("expected 201, got %d", res.Status)
				}

				commits, res, err := client.Git.ListCommits(context.Background(), repoFQ(testRepoName), scm.CommitListOptions{})
				if err != nil {
					t.Errorf("could not list commits: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if len(commits) != 2 {
					t.Errorf("expected 2 commits, got %d", len(commits))
				}
			},
		},
		{
			Name: "find commit",
			Test: func(t *testing.T) {
				commit, res, err := client.Git.FindCommit(context.Background(), repoFQ(testRepoName), readmeCommitSha)
				if err != nil {
					t.Errorf("could not find commit: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if commit.Author.Email != "amca@example.com" {
					t.Errorf("expected author email 'amca@example.com', got '%s'", commit.Author.Email)
				}

			},
		},
		{
			Name: "compare commits",
			Test: func(t *testing.T) {
				commits, _, err := client.Git.ListCommits(context.Background(), repoFQ(testRepoName), scm.CommitListOptions{})
				if err != nil {
					t.Errorf("could not list commits %v", err)
				}

				changes, res, err := client.Git.CompareCommits(context.Background(), repoFQ(testRepoName), commits[1].Sha, commits[0].Sha, &scm.ListOptions{})
				if err != nil {
					t.Errorf("could not compare commits: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if len(changes) == 0 {
					t.Error("there were no changes found")
				}
			},
		},
	}

	setup(t)

	for _, tc := range tests {
		if !t.Run(tc.Name, tc.Test) {
			t.Fatal("test aborted")
		}
	}
}
