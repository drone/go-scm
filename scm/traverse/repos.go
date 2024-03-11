// Copyright 2022 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traverse

import (
	"context"

	"github.com/drone/go-scm/scm"
	"golang.org/x/sync/errgroup"
)

// Repos returns the full repository list, traversing and
// combining paginated responses if necessary.
func Repos(ctx context.Context, client *scm.Client) ([]*scm.Repository, error) {
	list := []*scm.Repository{}
	opts := scm.ListOptions{Size: 100}
	for {
		result, meta, err := client.Repositories.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		list = addNonNil(list, result)
		opts.Page = meta.Page.Next
		opts.URL = meta.Page.NextURL

		if opts.Page == 0 && opts.URL == "" {
			break
		}
	}
	return list, nil
}

// ReposV2 same as Repos but uses errgroup to fetch repos in parallel
func ReposV2(ctx context.Context, client *scm.Client) ([]*scm.Repository, error) {
	list := []*scm.Repository{}
	opts := scm.ListOptions{Size: 100}

	result, meta, err := client.Repositories.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	list = addNonNil(list, result)
	if meta.Page.Next == 0 && meta.Page.NextURL == "" {
		return list, nil
	}
	errGroup, ectx := errgroup.WithContext(ctx)
	for i := meta.Page.Next; i <= meta.Page.Last; i++ {
		opts := scm.ListOptions{Size: 100, Page: i}
		errGroup.Go(func() error {
			result, _, err := client.Repositories.List(ectx, opts)
			if err != nil {
				return err
			}
			list = addNonNil(list, result)
			return nil
		})
	}
	return list, errGroup.Wait()
}

func addNonNil(list []*scm.Repository, result []*scm.Repository) []*scm.Repository {
	for _, src := range result {
		if src != nil {
			list = append(list, src)
		}
	}
	return list
}
