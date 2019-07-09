package scm

import (
	"fmt"
	"strings"
)

// MissingUsers is an error specifying the users that could not be unassigned.
type MissingUsers struct {
	Users  []string
	Action string
}

func (m MissingUsers) Error() string {
	return fmt.Sprintf("could not %s the following user(s): %s.", m.Action, strings.Join(m.Users, ", "))
}

// ExtraUsers is an error specifying the users that could not be unassigned.
type ExtraUsers struct {
	Users  []string
	Action string
}

func (e ExtraUsers) Error() string {
	return fmt.Sprintf("could not %s the following user(s): %s.", e.Action, strings.Join(e.Users, ", "))
}

// UnknownWebhook if the webhook is unknown
type UnknownWebhook struct {
	Event string
}

func (e UnknownWebhook) Error() string {
	return fmt.Sprintf("Unknown webhook event: %s.", e.Event)
}

// IsUnknownWebhook returns true if the error is an unknown webhook
func IsUnknownWebhook(err error) bool {
	_, ok := err.(UnknownWebhook)
	return ok
}
