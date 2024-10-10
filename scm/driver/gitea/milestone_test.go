package gitea

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestMilestoneFind(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://demo.gitea.com")
	got, _, err := client.Milestones.Find(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Milestone)
	raw, _ := os.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneList(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Get("/api/v1/repos/jcitizen/my-repo/milestones").
		Reply(200).
		Type("application/json").
		SetHeaders(mockPageHeaders).
		File("testdata/milestones.json")

	client, _ := New("https://demo.gitea.com")
	got, res, err := client.Milestones.List(context.Background(), "jcitizen/my-repo", scm.MilestoneListOptions{})
	if err != nil {
		t.Error(err)
	}

	want := []*scm.Milestone{}
	raw, _ := os.ReadFile("testdata/milestones.json.golden")
	err = json.Unmarshal(raw, &want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Page", testPage(res))
}

func TestMilestoneCreate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Post("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://demo.gitea.com")
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     &dueDate,
	}
	got, _, err := client.Milestones.Create(context.Background(), "jcitizen/my-repo", input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Milestone)
	raw, _ := os.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneUpdate(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Patch("/api/v1/repos/jcitizen/my-repo/milestones").
		File("testdata/milestone_create.json").
		Reply(200).
		Type("application/json").
		File("testdata/milestone.json")

	client, _ := New("https://demo.gitea.com")
	dueDate, _ := time.Parse(scm.SearchTimeFormat, "2012-10-09T23:39:01Z")
	input := &scm.MilestoneInput{
		Title:       "v1.0",
		Description: "Tracking milestone for version 1.0",
		State:       "open",
		DueDate:     &dueDate,
	}
	got, _, err := client.Milestones.Update(context.Background(), "jcitizen/my-repo", 1, input)
	if err != nil {
		t.Error(err)
	}

	want := new(scm.Milestone)
	raw, _ := os.ReadFile("testdata/milestone.json.golden")
	err = json.Unmarshal(raw, want)
	assert.NoError(t, err)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestMilestoneDelete(t *testing.T) {
	defer gock.Off()

	mockServerVersion()

	gock.New("https://demo.gitea.com").
		Delete("/api/v1/repos/jcitizen/my-repo/milestones/1").
		Reply(200).
		Type("application/json")

	client, _ := New("https://demo.gitea.com")
	_, err := client.Milestones.Delete(context.Background(), "jcitizen/my-repo", 1)
	if err != nil {
		t.Error(err)
	}
}
