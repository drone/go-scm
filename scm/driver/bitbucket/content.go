// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"bytes"
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type contentService struct {
	client *wrapper
}

func (s *contentService) Find(ctx context.Context, repo, path, ref string) (*scm.Content, *scm.Response, error) {
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src/%s/%s", repo, ref, path)
	out := new(bytes.Buffer)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return &scm.Content{
		Path: path,
		Data: out.Bytes(),
	}, res, err
}

func (s *contentService) Create(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) Update(ctx context.Context, repo, path string, params *scm.ContentParams) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) Delete(ctx context.Context, repo, path, ref string) (*scm.Response, error) {
	return nil, scm.ErrNotSupported
}

func (s *contentService) List(ctx context.Context, repo, path, ref string, opts scm.ListOptions) ([]*scm.ContentInfo, *scm.Response, error) {
	endpoint := fmt.Sprintf("/2.0/repositories/%s/src/%s/%s?%s", repo, ref, path, encodeListOptions(opts))
	out := new(contents)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	copyPagination(out.pagination, res)
	return convertContentInfoList(out), res, err
}

type contents struct {
	pagination
	Values []*content `json:"values"`
}

type content struct {
	Path       string   `json:"path"`
	Type       string   `json:"type"`
	Attributes []string `json:"attributes"`
}

func convertContentInfoList(from *contents) []*scm.ContentInfo {
	to := []*scm.ContentInfo{}
	for _, v := range from.Values {
		to = append(to, convertContentInfo(v))
	}
	return to
}

func convertContentInfo(from *content) *scm.ContentInfo {
	to := &scm.ContentInfo{Path: from.Path}
	switch from.Type {
	case "commit_file":
		to.Kind = func() scm.ContentKind {
			for _, attr := range from.Attributes {
				switch attr {
				case "link":
					return scm.ContentKindSymlink
				case "subrepository":
					return scm.ContentKindGitlink
				}
			}
			return scm.ContentKindFile
		}()
	case "commit_directory":
		to.Kind = scm.ContentKindDirectory
	default:
		to.Kind = scm.ContentKindUnsupported
	}
	return to
}
