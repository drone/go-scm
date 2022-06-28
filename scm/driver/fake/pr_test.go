package fake

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/jenkins-x/go-scm/scm"
)

func TestListChangesPagination(t *testing.T) {
	prNum := 11
	pageTests := []struct {
		prNum     int
		items     int
		page      int
		size      int
		wantFiles []string
	}{
		{prNum, 10, 2, 5, []string{"file6", "file7", "file8", "file9", "file10"}},
		{50, 10, 2, 5, []string{}},
	}

	for i, tt := range pageTests {
		t.Run(fmt.Sprintf("[%d]", i+1), func(rt *testing.T) {
			ctx := context.Background()
			client, data := NewDefault()
			// This stores the data in the "prNum" PR, but the list gets it from
			// the test number.
			data.PullRequestChanges[prNum] = makeChanges(tt.items)

			items, _, err := client.PullRequests.ListChanges(ctx, "test/test", tt.prNum, &scm.ListOptions{Page: tt.page, Size: tt.size})
			if err != nil {
				t.Error(err)
				return
			}
			if got := extractChangeFiles(items); !reflect.DeepEqual(got, tt.wantFiles) {
				rt.Errorf("ListChanges() got %#v, want %#v", got, tt.wantFiles)
			}
		})
	}
}

func TestPaginated(t *testing.T) {
	tests := []struct {
		page      int
		size      int
		items     int
		wantStart int
		wantEnd   int
	}{
		{1, 5, 10, 0, 5},
		{2, 5, 10, 5, 10},
		{2, 5, 9, 5, 9},
		{4, 5, 10, 10, 10}, // this results in an empty slice
		{0, 0, 10, 0, 10},  // this is the default 0 value for ListOption
	}

	for _, tt := range tests {
		start, end := paginated(tt.page, tt.size, tt.items)
		if tt.wantStart != start || tt.wantEnd != end {
			t.Fatalf("paginaged(%d, %d, %d) got items[%d:%d], want items[%d:%d]", tt.page, tt.size, tt.items, start, end, tt.wantStart, tt.wantEnd)
		}
	}
}

func makeChanges(n int) []*scm.Change {
	c := []*scm.Change{}
	for i := 1; i <= n; i++ {
		c = append(c, &scm.Change{
			Path: fmt.Sprintf("file%d", i),
		})
	}
	return c
}

func extractChangeFiles(ch []*scm.Change) []string {
	f := []string{}
	for _, c := range ch {
		f = append(f, c.Path)
	}
	return f
}

func TestPullServiceList(t *testing.T) {
	ctx := context.Background()
	client, data := NewDefault()

	A := &scm.PullRequest{
		Number: 0,
		Base: scm.PullRequestBranch{
			Repo: scm.Repository{
				FullName: "test/test",
			},
		},
		Closed:  false,
		Created: time.Unix(100, 0),
		Updated: time.Unix(200, 0),
	}

	B := &scm.PullRequest{
		Number: 1,
		Base: scm.PullRequestBranch{
			Repo: scm.Repository{
				FullName: "test/test",
			},
		},
		Closed:  true,
		Created: time.Unix(100, 0),
		Updated: time.Unix(200, 0),
	}

	C := &scm.PullRequest{
		Number: 1,
		Base: scm.PullRequestBranch{
			Repo: scm.Repository{
				FullName: "test/test",
			},
		},
		Closed:  false,
		Created: time.Unix(50, 0),
		Updated: time.Unix(300, 0),
	}

	D := &scm.PullRequest{
		Number: 1,
		Base: scm.PullRequestBranch{
			Repo: scm.Repository{
				FullName: "test/test",
			},
		},
		Closed:  true,
		Created: time.Unix(150, 0),
		Updated: time.Unix(175, 0),
	}

	data.PullRequests[0] = A
	data.PullRequests[1] = B
	data.PullRequests[2] = C
	data.PullRequests[3] = D

	standardCreationTime := time.Unix(100, 0)
	standardUpdateTime := time.Unix(200, 0)

	justBeforeStandardCreation := standardCreationTime.Add(time.Second * -1)
	justAfterStandardCreation := standardCreationTime.Add(time.Second * 1)
	justBeforeStandardUpdate := standardUpdateTime.Add(time.Second * -1)
	justAfterStandardUpdate := standardUpdateTime.Add(time.Second * 1)

	listTests := []struct {
		Page           int
		Size           int
		Open           bool
		Closed         bool
		Labels         []string
		UpdatedAfter   *time.Time
		UpdatedBefore  *time.Time
		CreatedAfter   *time.Time
		CreatedBefore  *time.Time
		ExpectedResult []*scm.PullRequest
		Description    string
	}{
		{0, 0, true, true, nil, nil, nil, nil, nil, []*scm.PullRequest{A, B, C, D}, "all prs"},
		{0, 0, true, false, nil, nil, nil, nil, nil, []*scm.PullRequest{A, C}, "open prs"},
		{0, 0, false, true, nil, nil, nil, nil, nil, []*scm.PullRequest{B, D}, "closed prs"},
		{0, 0, false, false, nil, nil, nil, nil, nil, []*scm.PullRequest{}, "neither open, nor closed"},
		{1, 2, true, true, nil, nil, nil, nil, nil, []*scm.PullRequest{A, B}, "first page of prs"},
		{2, 2, true, true, nil, nil, nil, nil, nil, []*scm.PullRequest{C, D}, "second page of prs"},
		{1, 5, true, true, nil, nil, nil, nil, nil, []*scm.PullRequest{A, B, C, D}, "overlarge page of prs"},
		{2, 3, true, true, nil, nil, nil, nil, nil, []*scm.PullRequest{D}, "page of prs that doesn't fill size"},
		{0, 0, true, true, nil, nil, nil, nil, &standardCreationTime, []*scm.PullRequest{C}, "early creation prs"},
		{0, 0, true, true, nil, nil, nil, &standardCreationTime, nil, []*scm.PullRequest{D}, "late creation prs"},
		{0, 0, true, true, nil, nil, &standardUpdateTime, nil, nil, []*scm.PullRequest{D}, "early update prs"},
		{0, 0, true, true, nil, &standardUpdateTime, nil, nil, nil, []*scm.PullRequest{C}, "late update prs"},
		{0, 0, true, true, nil, &justBeforeStandardUpdate, &justAfterStandardUpdate, &justBeforeStandardCreation, &justAfterStandardCreation, []*scm.PullRequest{A, B}, "standard time created and updated prs"},
	}

	for i, tt := range listTests {
		items, _, err := client.PullRequests.List(ctx, "test/test", &scm.PullRequestListOptions{
			Page:          tt.Page,
			Size:          tt.Size,
			Open:          tt.Open,
			Closed:        tt.Closed,
			Labels:        tt.Labels,
			UpdatedAfter:  tt.UpdatedAfter,
			UpdatedBefore: tt.UpdatedBefore,
			CreatedAfter:  tt.CreatedAfter,
			CreatedBefore: tt.CreatedBefore,
		})

		if err != nil {
			t.Error(err)
			return
		}

		if len(items) != len(tt.ExpectedResult) {
			t.Error(fmt.Errorf("PullRequests.List(), iteration #%d: %s, expected len of \n%#v \nto match len of\n%#v", i, tt.Description, tt.ExpectedResult, items))
		}

		for _, item := range items {
			if !contains(tt.ExpectedResult, item) {
				t.Error(fmt.Errorf("PullRequests.List(), iteration #%d: %s, \nexpected %#v \nto contain \n%#v", i, tt.Description, tt.ExpectedResult, item))
			}
		}
	}
}

func contains(s []*scm.PullRequest, e *scm.PullRequest) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}
