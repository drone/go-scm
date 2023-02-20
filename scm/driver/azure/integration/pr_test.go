package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

func TestPullRequests(t *testing.T) {
	canRun(t)

	const testRepoName = "pull-request-test"
	var testRepo *scm.Repository
	var readmeCommitSha string
	var pullRequest *scm.PullRequest

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
			Name: "create a pr",
			Test: func(t *testing.T) {
				content := scm.ContentParams{
					Message: "create a main.go app",
					Data:    []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello world!\")\n}"),
					Branch:  "refs/heads/hello-world",
					Ref:     readmeCommitSha,
				}

				_, err := client.Contents.Create(context.Background(), repoFQ(testRepoName), "/main.go", &content)
				if err != nil {
					t.Errorf("could not create main.go file and commit: %v", err)
				}

				prInput := &scm.PullRequestInput{
					Title: "introduce hello world app",
					Body:  "introduce hello world app as main.go",
					Head:  "hello-world",
					Base:  "main",
				}

				var res *scm.Response
				pullRequest, res, err = client.PullRequests.Create(context.Background(), repoFQ(testRepoName), prInput)
				if err != nil {
					t.Fatalf("failed to create pull request: %v", err)
				}

				if res.Status != http.StatusCreated {
					t.Fatalf("expected 201, got: %d", res.Status)
				}

				expectedLink := fmt.Sprintf("%s/pullrequest/%d", testRepo.Link, pullRequest.Number)
				if pullRequest.Link != expectedLink {
					t.Fatalf("expected link '%s', got '%s'", expectedLink, pullRequest.Link)
				}
			},
		},
		{
			Name: "find the pr by id",
			Test: func(t *testing.T) {
				_, res, err := client.PullRequests.Find(context.Background(), repoFQ(testRepoName), pullRequest.Number)
				if err != nil {
					t.Fatalf("could not get pull request: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got: %d", res.Status)
				}
			},
		},
		{
			Name: "list all prs",
			Test: func(t *testing.T) {
				prs, res, err := client.PullRequests.List(context.Background(), repoFQ(testRepoName), &scm.PullRequestListOptions{Closed: false})
				if err != nil {
					t.Fatalf("could not list pull requests: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got: %d", res.Status)
				}

				if len(prs) != 1 {
					t.Fatalf("expected 1 pr in list, got: %d", len(prs))
				}

				if prs[0].Number != pullRequest.Number {
					t.Fatalf("expected pr id %d, got: %d", pullRequest.Number, prs[0].Number)
				}
			},
		},
		{
			Name: "list commits on PR",
			Test: func(t *testing.T) {
				commits, res, err := client.PullRequests.ListCommits(context.Background(), repoFQ(testRepoName), pullRequest.Number, &scm.ListOptions{})
				if err != nil {
					t.Fatalf("could not list commits: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got: %d", res.Status)
				}

				if len(commits) != 1 {
					t.Fatalf("expected 1 commit in list, got: %d", len(commits))
				}

				if commits[0].Message != "create a main.go app" {
					t.Fatalf("commit message does not match, got: %s", commits[0].Message)
				}
			},
		},
		{
			Name: "close pr without merging",
			Test: func(t *testing.T) {
				res, err := client.PullRequests.Close(context.Background(), repoFQ(testRepoName), pullRequest.Number)
				if err != nil {
					t.Fatalf("could not close merge: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got: %d", res.Status)
				}

				prs, _, err := client.PullRequests.List(context.Background(), repoFQ(testRepoName), &scm.PullRequestListOptions{Closed: false})
				if err != nil {
					t.Fatalf("could not list pull requests: %v", err)
				}

				if len(prs) != 0 {
					t.Fatalf("expected 0 open prs, got: %d", len(prs))
				}
			},
		},
		{
			Name: "merge pr",
			Test: func(t *testing.T) {
				content := scm.ContentParams{
					Message: "create a main.go app to merge",
					Data:    []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello I Merged!\")\n}"),
					Branch:  "refs/heads/hello-merge",
					Ref:     readmeCommitSha,
				}

				_, err := client.Contents.Create(context.Background(), repoFQ(testRepoName), "/main.go", &content)
				if err != nil {
					t.Errorf("could not create main.go file and commit: %v", err)
				}

				prInput := &scm.PullRequestInput{
					Title: "introduce hello I merged app",
					Body:  "introduce hello I merged app as main.go",
					Head:  "hello-merge",
					Base:  "main",
				}

				var res *scm.Response
				pullRequest, res, err = client.PullRequests.Create(context.Background(), repoFQ(testRepoName), prInput)
				if err != nil {
					t.Fatalf("failed to create pull request: %v", err)
				}

				if res.Status != http.StatusCreated {
					t.Fatalf("expected 201, got: %d", res.Status)
				}

				res, err = client.PullRequests.Merge(context.Background(), repoFQ(testRepoName), pullRequest.Number, &scm.PullRequestMergeOptions{})

				retries := 0
				for err != nil && err.Error() == "patch accepted, but status still active" && retries < 3 {
					retries++
					t.Logf("azure did not accept the merge, retrying %d", retries)
					time.Sleep(time.Second * 2)
					res, err = client.PullRequests.Merge(context.Background(), repoFQ(testRepoName), pullRequest.Number, &scm.PullRequestMergeOptions{})
				}

				if err != nil {
					t.Fatalf("could not merge: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Fatalf("expected 200, got: %d", res.Status)
				}

				prs, _, err := client.PullRequests.List(context.Background(), repoFQ(testRepoName), &scm.PullRequestListOptions{Closed: false})
				if err != nil {
					t.Fatalf("could not list pull requests: %v", err)
				}

				if len(prs) != 0 {
					t.Fatalf("expected 0 open prs, got: %d", len(prs))
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
