// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm_test

import (
	"context"
	"log"
	"net/http"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
)

var ctx context.Context

func ExampleClient() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	// Sets a custom http.Client. This can be used with
	// github.com/golang/oauth2 for authorization.
	client.Client = &http.Client{}
}

func ExampleUser_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	user, _, err := client.Users.Find(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user.Login)
}

func ExampleUser_findLogin() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	user, _, err := client.Users.FindLogin(ctx, "octocat")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user.Login)
}

func ExampleRepository_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	repo, _, err := client.Repositories.Find(ctx, "octocat/Hello-World")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(repo.Namespace, repo.Name)
}

func ExampleRepository_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.ListOptions{
		Page: 1,
		Size: 30,
	}

	repos, _, err := client.Repositories.List(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		log.Println(repo.Namespace, repo.Name)
	}
}

func ExampleIssue_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.IssueListOptions{
		Page:   1,
		Size:   30,
		Open:   true,
		Closed: false,
	}

	issues, _, err := client.Issues.List(ctx, "octocat/Hello-World", opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issues {
		log.Println(issue.Number, issue.Title)
	}
}

func ExampleIssue_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	issue, _, err := client.Issues.Find(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(issue.Number, issue.Title)
}

func ExampleIssue_close() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Issues.Close(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleIssue_lock() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Issues.Lock(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleIssue_unlock() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Issues.Lock(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleOrganization_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.ListOptions{
		Page: 1,
		Size: 30,
	}

	orgs, _, err := client.Organizations.List(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, org := range orgs {
		log.Println(org.Name)
	}
}

func ExampleOrganization_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	org, _, err := client.Organizations.Find(ctx, "github")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(org.Name)
}

func ExampleComment_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.ListOptions{
		Page: 1,
		Size: 30,
	}

	comments, _, err := client.Issues.ListComments(ctx, "octocat/Hello-World", 1, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, comment := range comments {
		log.Println(comment.Body)
	}
}

func ExampleComment_create() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	in := &scm.CommentInput{
		Body: "Found a bug",
	}

	comment, _, err := client.Issues.CreateComment(ctx, "octocat/Hello-World", 1, in)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(comment.ID)
}

func ExampleComment_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	comment, _, err := client.Issues.FindComment(ctx, "octocat/Hello-World", 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(comment.Body)
}

func ExampleReview_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.ListOptions{
		Page: 1,
		Size: 30,
	}

	reviews, _, err := client.Reviews.List(ctx, "octocat/Hello-World", 1, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, review := range reviews {
		log.Println(
			review.Path,
			review.Line,
			review.Body,
		)
	}
}

func ExampleReview_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	review, _, err := client.Reviews.Find(ctx, "octocat/Hello-World", 1, 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(
		review.Path,
		review.Line,
		review.Body,
	)
}

func ExampleReview_create() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	in := &scm.ReviewInput{
		Line: 38,
		Path: "main.go",
		Body: "Run gofmt please",
	}

	review, _, err := client.Reviews.Create(ctx, "octocat/Hello-World", 1, in)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(review.ID)
}

func ExamplePullRequest_list() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	opts := scm.PullRequestListOptions{
		Page:   1,
		Size:   30,
		Open:   true,
		Closed: false,
	}

	prs, _, err := client.PullRequests.List(ctx, "octocat/Hello-World", opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, pr := range prs {
		log.Println(pr.Number, pr.Title)
	}
}

func ExamplePullRequest_find() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	pr, _, err := client.PullRequests.Find(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pr.Number, pr.Title)
}

func ExamplePullRequest_close() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.PullRequests.Close(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}
}

func ExamplePullRequest_merge() {
	client, err := github.New("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.PullRequests.Merge(ctx, "octocat/Hello-World", 1)
	if err != nil {
		log.Fatal(err)
	}
}
