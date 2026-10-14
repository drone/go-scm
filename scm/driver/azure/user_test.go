package azure

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

func TestUserFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com").
		Get("/ORG/_apis/connectionData").
		Reply(200).
		Type("application/json").
		File("testdata/connection_data.json")

	client := NewDefault("ORG", "")
	got, _, err := client.Users.Find(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	want := &scm.User{
		ID:    "020710a8-1111-2222-3333-444444444444",
		Name:  "AzureDevOpsUser2",
		Login: "AzureDevOpsUser2",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestUserFind_Error(t *testing.T) {
	defer gock.Off()

	gock.New("https://dev.azure.com").
		Get("/ORG/_apis/connectionData").
		Reply(401).
		Type("application/json").
		BodyString(`{"message":"unauthorized"}`)

	client := NewDefault("ORG", "")
	if _, _, err := client.Users.Find(context.Background()); err == nil {
		t.Error("expected error")
	}
}
