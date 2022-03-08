package gitlab

import (
	"context"
	"fmt"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateCommitStatus(t *testing.T) {
	defer gock.Off()
	sha := "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"

	gock.New("https://gitlab.com").
		Post(fmt.Sprintf("api/v4/projects/devops/demo/statuses/%s", sha)).
		MatchType("json").
		JSON(map[string]interface{}{
			"id":          "",
			"sha":         "",
			"ref":         "",
			"state":       "pending",
			"name":        "CodeScan",
			"description": "CodeScan Description",
			"pipeline_id": 29355,
			"target_url":  "https://gitlab.com",
			"coverage":    0,
		}).
		Reply(201).
		JSON(map[string]interface{}{
			"id":            54905,
			"sha":           sha,
			"ref":           "develop",
			"status":        "pending",
			"name":          "CodeScan",
			"target_url":    "https://gitlab.com",
			"description":   "",
			"created_at":    "2021-01-03T05:44:52.715Z",
			"started_at":    nil,
			"finished_at":   nil,
			"allow_failure": false,
			"coverage":      0,
			"author": map[string]interface{}{
				"id":         22,
				"name":       "Yin",
				"username":   "Yin",
				"state":      "active",
				"avatar_url": "",
				"web_url":    "",
			},
		}).
		Type("application/json").
		SetHeaders(mockHeaders)
	pipelineID := 29355
	options := &scm.CommitStatusUpdateOptions{
		State:       "pending",
		Name:        "CodeScan",
		Description: "CodeScan Description",
		PipelineID:  &pipelineID,
		TargetURL:   "https://gitlab.com",
	}
	client := NewDefault()

	status, _, err := client.Commits.UpdateCommitStatus(context.Background(), "devops/demo", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", options)
	if err != nil {
		t.Fatal(err)
	}

	if status.Sha != sha {
		t.Errorf("sha value should be %s", sha)
	}
	if status.Status != "pending" {
		t.Error("status value should be pending")
	}
}

func TestUpdateCommitStatusWithEmptyPipelineId(t *testing.T) {
	defer gock.Off()
	sha := "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"

	gock.New("https://gitlab.com").
		Post(fmt.Sprintf("api/v4/projects/devops/demo/statuses/%s", sha)).
		MatchType("json").
		JSON(map[string]interface{}{
			"id":          "",
			"sha":         "",
			"ref":         "",
			"state":       "pending",
			"name":        "CodeScan",
			"description": "CodeScan Description",
			"target_url":  "https://gitlab.com",
			"coverage":    0,
		}).
		Reply(201).
		JSON(map[string]interface{}{
			"id":            54905,
			"sha":           sha,
			"ref":           "develop",
			"status":        "pending",
			"name":          "CodeScan",
			"target_url":    "https://gitlab.com",
			"description":   "",
			"created_at":    "2021-01-03T05:44:52.715Z",
			"started_at":    nil,
			"finished_at":   nil,
			"allow_failure": false,
			"coverage":      0,
			"author": map[string]interface{}{
				"id":         22,
				"name":       "Yin",
				"username":   "Yin",
				"state":      "active",
				"avatar_url": "",
				"web_url":    "",
			},
		}).
		Type("application/json").
		SetHeaders(mockHeaders)
	options := &scm.CommitStatusUpdateOptions{
		State:       "pending",
		Name:        "CodeScan",
		Description: "CodeScan Description",
		TargetURL:   "https://gitlab.com",
	}
	client := NewDefault()

	status, _, err := client.Commits.UpdateCommitStatus(context.Background(), "devops/demo", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", options)
	if err != nil {
		t.Fatal(err)
	}

	if status.Sha != sha {
		t.Errorf("sha value should be %s", sha)
	}
	if status.Status != "pending" {
		t.Error("status value should be pending")
	}
}

func TestUpdateCommitStatusWithZeroPipelineId(t *testing.T) {
	defer gock.Off()
	sha := "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d"

	gock.New("https://gitlab.com").
		Post(fmt.Sprintf("api/v4/projects/devops/demo/statuses/%s", sha)).
		MatchType("json").
		JSON(map[string]interface{}{
			"id":          "",
			"sha":         "",
			"ref":         "",
			"state":       "pending",
			"name":        "CodeScan",
			"pipeline_id": 0,
			"description": "CodeScan Description",
			"target_url":  "https://gitlab.com",
			"coverage":    0,
		}).
		Reply(201).
		JSON(map[string]interface{}{
			"id":            54905,
			"sha":           sha,
			"ref":           "develop",
			"status":        "pending",
			"name":          "CodeScan",
			"target_url":    "https://gitlab.com",
			"description":   "",
			"created_at":    "2021-01-03T05:44:52.715Z",
			"started_at":    nil,
			"finished_at":   nil,
			"allow_failure": false,
			"coverage":      0,
			"author": map[string]interface{}{
				"id":         22,
				"name":       "Yin",
				"username":   "Yin",
				"state":      "active",
				"avatar_url": "",
				"web_url":    "",
			},
		}).
		Type("application/json").
		SetHeaders(mockHeaders)
	pipelineID := 0
	options := &scm.CommitStatusUpdateOptions{
		State:       "pending",
		Name:        "CodeScan",
		Description: "CodeScan Description",
		TargetURL:   "https://gitlab.com",
		PipelineID:  &pipelineID,
	}
	client := NewDefault()

	status, _, err := client.Commits.UpdateCommitStatus(context.Background(), "devops/demo", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d", options)
	if err != nil {
		t.Fatal(err)
	}

	if status.Sha != sha {
		t.Errorf("sha value should be %s", sha)
	}
	if status.Status != "pending" {
		t.Error("status value should be pending")
	}
}
