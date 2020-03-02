// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/drone/go-scm/scm"
	"time"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref)
	out := new(content)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	raw, _ := base64.StdEncoding.DecodeString(out.Content)
	return &scm.Content{
		Path: out.Path,
		Data: raw,
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := new(contentCreateUpdate)
	in.Committer.Name = params.Signature.Name
	in.Committer.Email = params.Signature.Email
	in.Author.Name = params.Signature.Name
	in.Author.Email = params.Signature.Email
	in.Message = params.Message
	in.Branch = params.Branch
	in.Content = params.Data
	in.SHA = params.SHA
	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s", repo, path)
	in := new(contentCreateUpdate)
	in.Committer.Name = params.Signature.Name
	in.Committer.Email = params.Signature.Email
	in.Author.Name = params.Signature.Name
	in.Author.Email = params.Signature.Email
	in.Message = params.Message
	in.Branch = params.Branch
	in.Content = params.Data
	in.SHA = params.SHA
	res, err := s.client.do(ctx, "PUT", endpoint, in, nil)
	return res, err
}

func (s *contentService) Delete(ctx context.Context, repo, path, ref string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, _ scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("repos/%s/contents/%s?ref=%s", repo, path, ref)
	out := []*content{}
	res, err := s.client.do(ctx, "GET", endpoint, nil, &out)
	return convertContentInfoList(out), res, err
}

type content struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Sha     string `json:"sha"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type contentCreateUpdate struct {
	Branch    string
	Message   string
	Content   []byte
	SHA       string
	Author    CommitAuthor
	Committer CommitAuthor
}

type CommitAuthor struct {
	Date  *time.Time
	Name  string
	Email string
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
	switch from.Type {
	case "file":
		to.Kind = scm.ContentKindFile
	case "dir":
		to.Kind = scm.ContentKindDirectory
	case "symlink":
		to.Kind = scm.ContentKindSymlink
	case "submodule":
		to.Kind = scm.ContentKindGitlink
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
