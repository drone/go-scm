// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scm

import (
	"encoding/json"
)

// State represents the commit state.
type State int

// State values.
const (
	StateUnknown State = iota
	StatePending
	StateRunning
	StateSuccess
	StateFailure
	StateCanceled
	StateError
)

// String returns a string representation of the State
func (s State) String() string {
	switch s {
	case StateUnknown:
		return "unknown"
	case StatePending:
		return "pending"
	case StateRunning:
		return "running"
	case StateSuccess:
		return "success"
	case StateFailure:
		return "failure"
	case StateCanceled:
		return "cancelled"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// Action identifies webhook actions.
type Action int

// Action values.
const (
	ActionCreate Action = iota + 1
	ActionUpdate
	ActionDelete
	// issues
	ActionOpen
	ActionReopen
	ActionClose
	ActionLabel
	ActionUnlabel
	// pull requests
	ActionSync
	ActionMerge
	ActionAssigned
	ActionUnassigned
	ActionReviewRequested
	ActionReviewRequestRemoved
	ActionReadyForReview

	// reviews
	ActionEdited
	ActionSubmitted
	ActionDismissed
)

// String returns the string representation of Action.
func (a Action) String() (s string) {
	switch a {
	case ActionCreate:
		return "created"
	case ActionUpdate:
		return "updated"
	case ActionDelete:
		return "deleted"
	case ActionLabel:
		return "labeled"
	case ActionUnlabel:
		return "unlabeled"
	case ActionOpen:
		return "opened"
	case ActionReopen:
		return "reopened"
	case ActionClose:
		return "closed"
	case ActionSync:
		return "synchronized"
	case ActionMerge:
		return "merged"
	case ActionEdited:
		return "edited"
	case ActionSubmitted:
		return "submitted"
	case ActionDismissed:
		return "dismisssed"
	case ActionAssigned:
		return "assigned"
	case ActionUnassigned:
		return "unassigned"
	case ActionReviewRequested:
		return "review_requested"
	case ActionReviewRequestRemoved:
		return "review_request_removed"
	case ActionReadyForReview:
		return "ready_for_review"
	default:
		return
	}
}

// MarshalJSON returns the JSON-encoded Action.
func (a Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON unmarshales the JSON-encoded Action.
func (a *Action) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "created":
		*a = ActionCreate
	case "updated":
		*a = ActionUpdate
	case "deleted":
		*a = ActionDelete
	case "labeled":
		*a = ActionLabel
	case "unlabeled":
		*a = ActionUnlabel
	case "opened":
		*a = ActionOpen
	case "reopened":
		*a = ActionReopen
	case "closed":
		*a = ActionClose
	case "synchronized":
		*a = ActionSync
	case "merged":
		*a = ActionMerge
	}
	return nil
}

// Driver identifies source code management driver.
type Driver int

// Driver values.
const (
	DriverUnknown Driver = iota
	DriverGithub
	DriverGitlab
	DriverGogs
	DriverGitea
	DriverBitbucket
	DriverStash
	DriverCoding
	DriverFake
)

// String returns the string representation of Driver.
func (d Driver) String() (s string) {
	switch d {
	case DriverGithub:
		return "github"
	case DriverGitlab:
		return "gitlab"
	case DriverGogs:
		return "gogs"
	case DriverGitea:
		return "gitea"
	case DriverBitbucket:
		return "bitbucket"
	case DriverStash:
		return "stash"
	case DriverCoding:
		return "coding"
	case DriverFake:
		return "fake"
	default:
		return "unknown"
	}
}
