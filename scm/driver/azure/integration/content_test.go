package integration_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

func TestContentManagement(t *testing.T) {
	canRun(t)

	const testRepoName = "content-management-test"
	var testRepo *scm.Repository
	var readmeCommitSha string

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
			Name: "find README.md",
			Test: func(t *testing.T) {
				content, res, err := client.Contents.Find(context.Background(), repoFQ(testRepoName), "README.md", "")
				if err != nil {
					t.Errorf("could not find the README.md file: %v", err)
				}

				if string(content.Data) != "# Content Test" {
					t.Errorf("expected README heading, got: %s", content.Data)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if content.Sha != readmeCommitSha {
					t.Errorf("expected a sha %s, got %s", readmeCommitSha, content.Sha)
				}
			},
		},
		{
			Name: "create a main.go file",
			Test: func(t *testing.T) {
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
			},
		},
		{
			Name: "list content",
			Test: func(t *testing.T) {
				files, res, err := client.Contents.List(context.Background(), repoFQ(testRepoName), "", "", &scm.ListOptions{})
				if err != nil {
					t.Errorf("could not list content: %v", err)
				}

				if res.Status != http.StatusOK {
					t.Errorf("expected 200, got %d", res.Status)
				}

				if len(files) != 3 {
					t.Errorf("expected 3 file entries, got %d", len(files))
				}
			},
		},
		{
			Name: "update the main.go file",
			Test: func(t *testing.T) {
				latestCommit, _, err := client.Contents.List(context.Background(), repoFQ(testRepoName), "", "main", &scm.ListOptions{})
				if err != nil {
					t.Errorf("could not get the latest commit: %v", err)
				}

				latestSha := latestCommit[0].Sha

				content := scm.ContentParams{
					Message: "modify hello world to hello you",
					Data:    []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello you!\")\n}"),
					Branch:  "refs/heads/main",
					Sha:     latestSha,
				}
				res, err := client.Contents.Update(context.Background(), repoFQ(testRepoName), "/main.go", &content)

				if err != nil {
					t.Errorf("could not update main.go: %v", err)
				}

				if res.Status != http.StatusCreated {
					t.Errorf("expected 201, got %d", res.Status)
				}
			},
		},
		{
			Name: "delete the main.go file",
			Test: func(t *testing.T) {
				latestCommit, _, err := client.Contents.List(context.Background(), repoFQ(testRepoName), "", "main", &scm.ListOptions{})
				if err != nil {
					t.Errorf("could not get the latest commit: %v", err)
				}

				latestSha := latestCommit[0].Sha

				content := scm.ContentParams{
					Message: "delete main.go",
					Branch:  "refs/heads/main",
					Sha:     latestSha,
				}
				res, err := client.Contents.Delete(context.Background(), repoFQ(testRepoName), "/main.go", &content)

				if err != nil {
					t.Errorf("could not delete main.go: %v", err)
				}

				if res.Status != http.StatusCreated {
					t.Errorf("expected 201, got %d", res.Status)
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
