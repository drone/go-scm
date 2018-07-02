// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bitbucket

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/go-scm/scm"
)

// TODO(bradrydzewski) default repository branch is missing in webhook payloads

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (interface{}, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}

	var hook interface{}
	switch req.Header.Get("x-event-key") {
	case "repo:push":
		hook, err = s.parsePushHook(data)
		// case "create":
		// 	hook, err = s.parseCreateHook(data)
		// case "delete":
		// 	hook, err = s.parseDeleteHook(data)
		// case "pull_request":
		// 	hook, err = s.parsePullRequestHook(data)
		// case "pull_request_review_comment":
		// case "issues":
		// case "issue_comment":
	}
	if err != nil {
		return nil, err
	}
	if hook == nil {
		return nil, nil
	}

	// get the gogs signature key to verify the payload
	// signature. If no key is provided, no validation
	// is performed.
	key, err := fn(hook)
	if err != nil {
		return hook, err
	} else if key == "" {
		return hook, nil
	}

	if req.FormValue("secret") != key {
		return hook, scm.ErrSignatureInvalid
	}

	return hook, nil
}

func (s *webhookService) parsePushHook(data []byte) (interface{}, error) {
	dst := new(pushHook)
	err := json.Unmarshal(data, dst)
	if err != nil {
		return nil, err
	}
	if len(dst.Push.Changes) == 0 {
		return nil, errors.New("Push hook has empty changeset")
	}
	change := dst.Push.Changes[0]
	switch {
	case change.New.Type == "branch" && change.Created:
		return convertBranchCreateHook(dst), nil
	case change.Old.Type == "branch" && change.Closed:
		return convertBranchDeleteHook(dst), nil
	case change.New.Type == "tag" && change.Created:
		return convertTagCreateHook(dst), nil
	case change.Old.Type == "tag" && change.Closed:
		return convertTagDeleteHook(dst), nil
	default:
		return convertPushHook(dst), err
	}
}

//
// native data structures
//

type (
	// 	// github create webhook payload
	// 	createDeleteHook struct {
	// 		Ref        string     `json:"ref"`
	// 		RefType    string     `json:"ref_type"`
	// 		Repository repository `json:"repository"`
	// 		Sender     user       `json:"sender"`
	// 	}

	// bitbucket push webhook payload
	pushHook struct {
		Push struct {
			Changes []struct {
				Forced bool `json:"forced"`
				Old    struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Links struct {
						Commits struct {
							Href string `json:"href"`
						} `json:"commits"`
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Target struct {
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
						Author struct {
							Raw  string `json:"raw"`
							Type string `json:"type"`
							User struct {
								Username    string `json:"username"`
								DisplayName string `json:"display_name"`
								AccountID   string `json:"account_id"`
								Links       struct {
									Self struct {
										Href string `json:"href"`
									} `json:"self"`
									HTML struct {
										Href string `json:"href"`
									} `json:"html"`
									Avatar struct {
										Href string `json:"href"`
									} `json:"avatar"`
								} `json:"links"`
								Type string `json:"type"`
								UUID string `json:"uuid"`
							} `json:"user"`
						} `json:"author"`
						Summary struct {
							Raw    string `json:"raw"`
							Markup string `json:"markup"`
							HTML   string `json:"html"`
							Type   string `json:"type"`
						} `json:"summary"`
						Parents []interface{} `json:"parents"`
						Date    time.Time     `json:"date"`
						Message string        `json:"message"`
						Type    string        `json:"type"`
					} `json:"target"`
				} `json:"old"`
				Links struct {
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
					Diff struct {
						Href string `json:"href"`
					} `json:"diff"`
				} `json:"links"`
				Truncated bool `json:"truncated"`
				Commits   []struct {
					Hash  string `json:"hash"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						Comments struct {
							Href string `json:"href"`
						} `json:"comments"`
						Patch struct {
							Href string `json:"href"`
						} `json:"patch"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
						Diff struct {
							Href string `json:"href"`
						} `json:"diff"`
						Approve struct {
							Href string `json:"href"`
						} `json:"approve"`
						Statuses struct {
							Href string `json:"href"`
						} `json:"statuses"`
					} `json:"links"`
					Author struct {
						Raw  string `json:"raw"`
						Type string `json:"type"`
						User struct {
							Username    string `json:"username"`
							DisplayName string `json:"display_name"`
							AccountID   string `json:"account_id"`
							Links       struct {
								Self struct {
									Href string `json:"href"`
								} `json:"self"`
								HTML struct {
									Href string `json:"href"`
								} `json:"html"`
								Avatar struct {
									Href string `json:"href"`
								} `json:"avatar"`
							} `json:"links"`
							Type string `json:"type"`
							UUID string `json:"uuid"`
						} `json:"user"`
					} `json:"author"`
					Summary struct {
						Raw    string `json:"raw"`
						Markup string `json:"markup"`
						HTML   string `json:"html"`
						Type   string `json:"type"`
					} `json:"summary"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Date    time.Time `json:"date"`
					Message string    `json:"message"`
					Type    string    `json:"type"`
				} `json:"commits"`
				Created bool `json:"created"`
				Closed  bool `json:"closed"`
				New     struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Links struct {
						Commits struct {
							Href string `json:"href"`
						} `json:"commits"`
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Target struct {
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
						Author struct {
							Raw  string `json:"raw"`
							Type string `json:"type"`
							User struct {
								Username    string `json:"username"`
								DisplayName string `json:"display_name"`
								AccountID   string `json:"account_id"`
								Links       struct {
									Self struct {
										Href string `json:"href"`
									} `json:"self"`
									HTML struct {
										Href string `json:"href"`
									} `json:"html"`
									Avatar struct {
										Href string `json:"href"`
									} `json:"avatar"`
								} `json:"links"`
								Type string `json:"type"`
								UUID string `json:"uuid"`
							} `json:"user"`
						} `json:"author"`
						Summary struct {
							Raw    string `json:"raw"`
							Markup string `json:"markup"`
							HTML   string `json:"html"`
							Type   string `json:"type"`
						} `json:"summary"`
						Parents []struct {
							Type  string `json:"type"`
							Hash  string `json:"hash"`
							Links struct {
								Self struct {
									Href string `json:"href"`
								} `json:"self"`
								HTML struct {
									Href string `json:"href"`
								} `json:"html"`
							} `json:"links"`
						} `json:"parents"`
						Date    time.Time `json:"date"`
						Message string    `json:"message"`
						Type    string    `json:"type"`
					} `json:"target"`
				} `json:"new"`
			} `json:"changes"`
		} `json:"push"`
		Repository webhookRepository `json:"repository"`
		Actor      webhookActor      `json:"actor"`
	}

	webhookRepository struct {
		Scm   string `json:"scm"`
		Name  string `json:"name"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
		FullName string `json:"full_name"`
		Owner    struct {
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			AccountID   string `json:"account_id"`
			Links       struct {
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
			UUID string `json:"uuid"`
		} `json:"owner"`
		IsPrivate bool   `json:"is_private"`
		UUID      string `json:"uuid"`
	}

	webhookActor struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		AccountID   string `json:"account_id"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		UUID string `json:"uuid"`
	}
)

//
// native data structure conversion
//

func convertPushHook(src *pushHook) *scm.PushHook {
	change := src.Push.Changes[0]
	return &scm.PushHook{
		Ref: "refs/heads/" + change.New.Name,
		Commit: scm.Commit{
			Sha:     change.New.Target.Hash,
			Message: change.New.Target.Message,
			Link:    change.New.Target.Links.HTML.Href,
			Author: scm.Signature{
				Login:  change.New.Target.Author.User.Username,
				Email:  extractEmail(change.New.Target.Author.Raw),
				Name:   change.New.Target.Author.User.DisplayName,
				Avatar: change.New.Target.Author.User.Links.Avatar.Href,
				Date:   change.New.Target.Date,
			},
			Committer: scm.Signature{
				Login:  change.New.Target.Author.User.Username,
				Email:  extractEmail(change.New.Target.Author.Raw),
				Name:   change.New.Target.Author.User.DisplayName,
				Avatar: change.New.Target.Author.User.Links.Avatar.Href,
				Date:   change.New.Target.Date,
			},
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: src.Repository.Owner.Username,
			Name:      src.Repository.Name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertBranchCreateHook(src *pushHook) *scm.BranchHook {
	change := src.Push.Changes[0].New
	action := scm.ActionCreate
	return &scm.BranchHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: src.Repository.Owner.Username,
			Name:      src.Repository.Name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertBranchDeleteHook(src *pushHook) *scm.BranchHook {
	change := src.Push.Changes[0].Old
	action := scm.ActionDelete
	return &scm.BranchHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: src.Repository.Owner.Username,
			Name:      src.Repository.Name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertTagCreateHook(src *pushHook) *scm.TagHook {
	change := src.Push.Changes[0].New
	action := scm.ActionCreate
	return &scm.TagHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: src.Repository.Owner.Username,
			Name:      src.Repository.Name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}

func convertTagDeleteHook(src *pushHook) *scm.TagHook {
	change := src.Push.Changes[0].Old
	action := scm.ActionDelete
	return &scm.TagHook{
		Action: action,
		Ref: scm.Reference{
			Name: change.Name,
			Sha:  change.Target.Hash,
		},
		Repo: scm.Repository{
			ID:        src.Repository.UUID,
			Namespace: src.Repository.Owner.Username,
			Name:      src.Repository.Name,
			Private:   src.Repository.IsPrivate,
			Clone:     fmt.Sprintf("https://bitbucket.org/%s.git", src.Repository.FullName),
			CloneSSH:  fmt.Sprintf("git@bitbucket.org:%s.git", src.Repository.FullName),
			Link:      src.Repository.Links.HTML.Href,
		},
		Sender: scm.User{
			Login:  src.Actor.Username,
			Name:   src.Actor.DisplayName,
			Avatar: src.Actor.Links.Avatar.Href,
		},
	}
}
