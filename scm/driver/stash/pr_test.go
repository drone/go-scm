// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stash

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.Find(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullFindComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.FindComment(context.Background(), "PRJ/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests").
		Reply(200).
		Type("application/json").
		File("testdata/prs.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.List(context.Background(), "PRJ/my-repo", scm.PullRequestListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/prs.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/changes").
		// MatchParam("pagelen", "30").
		// MatchParam("page", "1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_change.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.ListChanges(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_change.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/merge").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Merge(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/decline").
		Reply(200).
		Type("application/json").
		File("testdata/pr.json")

	client, _ := New("http://example.com:7990")
	_, err := client.PullRequests.Close(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests").
		Reply(201).
		Type("application/json").
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "Updated Files",
		Body:   `* added LICENSE\r\n* update files\r\n* update files`,
		Source: "feature/x",
		Target: "master",
	}

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.Create(context.Background(), "PRJ/my-repo", &input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPullCreateComment(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Post("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.CreateComment(context.Background(), "PRJ/my-repo", 1, &scm.CommentInput{
		Body: "LGTM",
	})
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

// TestEpochOrISO_UnmarshalJSON verifies that epochOrISO correctly accepts
// both integer epoch-milliseconds (pre-DC-10.3) and ISO 8601 strings (DC 10.3+).
func TestEpochOrISO_UnmarshalJSON(t *testing.T) {
	const wantMs = int64(1530766870000) // truncated to second precision

	tests := []struct {
		name    string
		input   string
		wantMs  int64
		wantErr bool
	}{
		{
			name:   "integer epoch-ms (pre-DC-10.3)",
			input:  `1530766870981`,
			wantMs: 1530766870981,
		},
		{
			name:   "ISO 8601 with +0000 offset no colon (DC 10.3+)",
			input:  `"2018-07-05T05:01:10+0000"`,
			wantMs: wantMs,
		},
		{
			name:   "ISO 8601 with +0200 offset no colon",
			input:  `"2018-07-05T07:01:10+0200"`,
			wantMs: wantMs,
		},
		{
			name:   "ISO 8601 RFC3339 with Z",
			input:  `"2018-07-05T05:01:10Z"`,
			wantMs: wantMs,
		},
		{
			name:   "ISO 8601 RFC3339 with colon offset",
			input:  `"2018-07-05T05:01:10+00:00"`,
			wantMs: wantMs,
		},
		{
			name:   "ISO 8601 with milliseconds and Z (DC 10.3 commit timestamps)",
			input:  `"2018-07-05T05:01:10.000Z"`,
			wantMs: wantMs,
		},
		{
			name:   "ISO 8601 with milliseconds and no-colon offset",
			input:  `"2018-07-05T05:01:10.000+0000"`,
			wantMs: wantMs,
		},
		{
			name:   "zero integer",
			input:  `0`,
			wantMs: 0,
		},
		{
			name:    "invalid string",
			input:   `"not-a-date"`,
			wantErr: true,
		},
		{
			name:   "null json (treated as zero by encoding/json)",
			input:  `null`,
			wantMs: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var v epochOrISO
			err := json.Unmarshal([]byte(tc.input), &v)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if int64(v) != tc.wantMs {
				t.Errorf("got %d ms, want %d ms (diff: %d ms)", int64(v), tc.wantMs, int64(v)-tc.wantMs)
			}
		})
	}
}

// TestEpochOrISO_TimeConversion verifies that the stored epoch-ms value
// produces the correct time.Time after /1000 (the call-site convention).
func TestEpochOrISO_TimeConversion(t *testing.T) {
	// epoch-ms for 2018-07-05T19:21:30+0000
	const epochMs = int64(1530818490000)
	want := time.Unix(epochMs/1000, 0)

	var v epochOrISO
	if err := json.Unmarshal([]byte(`"2018-07-05T19:21:30+0000"`), &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := time.Unix(int64(v)/1000, 0)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPullFindIsoDates(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_iso_dates.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.Find(context.Background(), "PRJ/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr_iso_dates.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ISO dates PR Find: unexpected results")
		t.Log(diff)
	}
}

func TestPullFindCommentIsoDates(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/comments/1").
		Reply(200).
		Type("application/json").
		File("testdata/pr_comment_iso_dates.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.FindComment(context.Background(), "PRJ/my-repo", 1, 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/pr_comment_iso_dates.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("ISO dates PR FindComment: unexpected results")
		t.Log(diff)
	}
}

func TestPullListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("http://example.com:7990").
		Get("rest/api/1.0/projects/PRJ/repos/my-repo/pull-requests/1/commits").
		Reply(200).
		Type("application/json").
		File("testdata/commits.json")

	client, _ := New("http://example.com:7990")
	got, _, err := client.PullRequests.ListCommits(context.Background(), "PRJ/my-repo", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	_ = json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}
