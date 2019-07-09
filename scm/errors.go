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
