// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package azure

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("%s/_apis/git/repositories/%s/items?path=%s&includeContent=true&$format=json&api-version=6.0", repo, s.client.AzureRepoID, path)
	out := new(content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return &scm.Content{
		Path:   out.Path,
		Data:   []byte(out.Content),
		Sha:    out.CommitID,
		BlobID: out.ObjectID,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}

	res, err := s.client.do(ctx, "POST", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		Content: params.Data,
		Sha:     params.Sha,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}
	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := &contentCreateUpdate{
		Message: params.Message,
		Branch:  params.Branch,
		Sha:     params.Sha,
		Committer: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
		Author: commitAuthor{
			Name:  params.Signature.Name,
			Email: params.Signature.Email,
		},
	}
	res, err := s.client.do(ctx, "DELETE", endpoint, in, nil)
	return res, err
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref)
	out := []*content{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	ObjectID      string `json:"objectId"`
	GitObjectType string `json:"gitObjectType"`
	CommitID      string `json:"commitId"`
	Path          string `json:"path"`
	Content       string `json:"content"`
	URL           string `json:"url"`
	Links         struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Repository struct {
			Href string `json:"href"`
		} `json:"repository"`
		Blob struct {
			Href string `json:"href"`
		} `json:"blob"`
	} `json:"_links"`
}

type contentCreateUpdate struct {
	Branch    string       `json:"branch"`
	Message   string       `json:"message"`
	Content   []byte       `json:"content"`
	Sha       string       `json:"sha"`
	Author    commitAuthor `json:"author"`
	Committer commitAuthor `json:"committer"`
}

type commitAuthor struct {
	Name  string `json:"name"`
	Date  string `json:"date"`
	Email string `json:"email"`
}

func convertContentInfoList(from []*content) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from {
		to = append(to, convertContentInfo(v))
	}
	return to
}
func convertContentInfo(from *content) *scm.ContentInfo {
	to := &scm.ContentInfo{Path: from.Path}
	// switch from.Type {
	// case "file":
	// 	to.Kind = scm.ContentKindFile
	// case "dir":
	// 	to.Kind = scm.ContentKindDirectory
	// case "symlink":
	// 	to.Kind = scm.ContentKindSymlink
	// case "submodule":
	// 	to.Kind = scm.ContentKindGitlink
	// default:
	// 	to.Kind = scm.ContentKindUnsupported
	// }
	return to
}
