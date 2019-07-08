package fake

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"k8s.io/apimachinery/pkg/util/sets"
)

type issueService struct {
	client *wrapper
	data   *Data
}

func (s *issueService) FindLabels(ctx context.Context, repo string, number int) ([]*scm.Label, *scm.Response, error) {
	f := s.data
	re := regexp.MustCompile(fmt.Sprintf(`^%s#%d:(.*)$`, repo, number))
	la := []*scm.Label{}
	allLabels := sets.NewString(f.IssueLabelsExisting...)
	allLabels.Insert(f.IssueLabelsAdded...)
	allLabels.Delete(f.IssueLabelsRemoved...)
	for _, l := range allLabels.List() {
		groups := re.FindStringSubmatch(l)
		if groups != nil {
			la = append(la, &scm.Label{Name: groups[1]})
		}
	}
	return la, nil, nil
}

func (s *issueService) Find(ctx context.Context, repo string, number int) (*scm.Issue, *scm.Response, error) {
	f := s.data
	for _, slice := range f.Issues {
		for _, issue := range slice {
			if issue.Number == number {
				return issue, nil, nil
			}
		}
	}
	return nil, nil, nil
}

func (s *issueService) AddLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	f := s.data
	labelString := fmt.Sprintf("%s#%d:%s", repo, number, label)
	if sets.NewString(f.IssueLabelsAdded...).Has(labelString) {
		return nil, fmt.Errorf("cannot add %v to %s/#%d", label, repo, number)
	}
	if f.RepoLabelsExisting == nil {
		f.IssueLabelsAdded = append(f.IssueLabelsAdded, labelString)
		return nil, nil
	}
	for _, l := range f.RepoLabelsExisting {
		if label == l {
			f.IssueLabelsAdded = append(f.IssueLabelsAdded, labelString)
			return nil, nil
		}
	}
	return nil, fmt.Errorf("cannot add %v to %s/#%d", label, repo, number)
}

// DeleteLabel removes a label
func (s *issueService) DeleteLabel(ctx context.Context, repo string, number int, label string) (*scm.Response, error) {
	f := s.data
	labelString := fmt.Sprintf("%s#%d:%s", repo, number, label)
	if !sets.NewString(f.IssueLabelsRemoved...).Has(labelString) {
		f.IssueLabelsRemoved = append(f.IssueLabelsRemoved, labelString)
		return nil, nil
	}
	return nil, fmt.Errorf("cannot remove %v from %s/#%d", label, repo, number)
}

// FindIssues returns f.Issues
func (s *issueService) FindIssues(query, sort string, asc bool) ([]scm.Issue, error) {
	f := s.data
	var issues []scm.Issue
	for _, slice := range f.Issues {
		for _, issue := range slice {
			issues = append(issues, *issue)
		}
	}
	return issues, nil
}

// AssignIssue adds assignees.
func (s *issueService) AssignIssue(owner, repo string, number int, assignees []string) error {
	f := s.data
	var m MissingUsers
	for _, a := range assignees {
		if a == "not-in-the-org" {
			m.Users = append(m.Users, a)
			continue
		}
		f.AssigneesAdded = append(f.AssigneesAdded, fmt.Sprintf("%s/%s#%d:%s", owner, repo, number, a))
	}
	if m.Users == nil {
		return nil
	}
	return m
}

// MissingUsers is an error specifying the users that could not be unassigned.
type MissingUsers struct {
	Users  []string
	action string
}

func (m MissingUsers) Error() string {
	return fmt.Sprintf("could not %s the following user(s): %s.", m.action, strings.Join(m.Users, ", "))
}

func (s *issueService) FindComment(context.Context, string, int, int) (*scm.Comment, *scm.Response, error) {
	panic("implement me")
}

func (s *issueService) List(context.Context, string, scm.IssueListOptions) ([]*scm.Issue, *scm.Response, error) {
	panic("implement me")
}

func (s *issueService) ListComments(context.Context, string, int, scm.ListOptions) ([]*scm.Comment, *scm.Response, error) {
	panic("implement me")
}

func (s *issueService) Create(context.Context, string, *scm.IssueInput) (*scm.Issue, *scm.Response, error) {
	panic("implement me")
}

func (s *issueService) CreateComment(context.Context, string, int, *scm.CommentInput) (*scm.Comment, *scm.Response, error) {
	panic("implement me")
}

func (s *issueService) DeleteComment(context.Context, string, int, int) (*scm.Response, error) {
	panic("implement me")
}

func (s *issueService) Close(context.Context, string, int) (*scm.Response, error) {
	panic("implement me")
}

func (s *issueService) Lock(context.Context, string, int) (*scm.Response, error) {
	panic("implement me")
}

func (s *issueService) Unlock(context.Context, string, int) (*scm.Response, error) {
	panic("implement me")
}
