// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"testing"

	"github.com/drone/go-scm/scm"
	"github.com/h2non/gock"
)

func TestClient_Base(t *testing.T) {
	client, err := New("https://dev.azure.com", "org", "proj")
	if err != nil {
		t.Error(err)
	}
	got, want := client.BaseURL.String(), "https://dev.azure.com/"
	if got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Default(t *testing.T) {
	client := NewDefault("org", "proj")
	if got, want := client.BaseURL.String(), "https://dev.azure.com/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_azure_special(t *testing.T) {
	client, _ := New("https://dev.azure.com", "org", "")
	if got, want := client.BaseURL.String(), "https://dev.azure.com/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
	client2, _ := New("https://dev.azure.com", "", "proj")
	if client2 != nil {
		t.Errorf("Want nil client, got %v", client2)
	}
}

func TestSanitizeBranchName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"master",
			args{
				"master",
			},
			"refs/heads/master",
		},
		{
			"feature/main-patch",
			args{
				"feature/main-patch",
			},
			"refs/heads/feature/main-patch",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeBranchName(tt.args.name); got != tt.want {
				t.Errorf("SanitizeBranchName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_AzureDevOpsServer(t *testing.T) {
	defer gock.Off()

	gock.New("https://host/tfs/").
		Get("/DefaultCollection/Proj/_apis/git/repositories/REPOID/refs").
		MatchParam("api-version", "^6.0$").
		Reply(200).
		Type("application/json").
		File("testdata/tags.json")

	client, err := New("https://host/tfs/", "DefaultCollection", "Proj")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := client.BaseURL.String(), "https://host/tfs/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
	if _, _, err := client.Git.ListTags(context.Background(), "REPOID", scm.ListOptions{}); err != nil {
		t.Error(err)
	}
	if !gock.IsDone() {
		t.Error("expected request to https://host/tfs/DefaultCollection/Proj/_apis/... was not made")
	}
}
