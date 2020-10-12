package fake_test

import (
	"context"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInviteLogic(t *testing.T) {
	client, data := fake.NewDefault()

	ctx := context.Background()
	owner := "myorg"
	repo := "myrepo"
	permission := "admin"
	user := "jstrachan"
	fullName := scm.Join(owner, repo)

	invitations, _, err := client.Users.ListInvitations(ctx)
	require.NoError(t, err, "could not list invitations in repo %s", fullName)
	require.Empty(t, invitations, "should not have any invitations")

	addedFlag, alreadyExisted, _, err := client.Repositories.AddCollaborator(ctx, fullName, user, permission)
	assert.True(t, addedFlag, "should have added a collaborator")
	assert.False(t, alreadyExisted, "the collaborator %s should not already exist", user)

	require.NotEmpty(t, data.UserPermissions[fullName], "should have a user permission for repo %s", fullName)
	assert.Equal(t, permission, data.UserPermissions[fullName][user], "should have a permission for repo %s and user %s", fullName, user)

	invitations, _, err = client.Users.ListInvitations(ctx)
	require.NoError(t, err, "could not list invitations in repo %s", fullName)
	require.Len(t, invitations, 1, "should have one invitation")

	invitation := invitations[0]
	t.Logf("we now have an invite %v for repo %s\n", invitation.ID, invitation.Repo.FullName)

	_, err = client.Users.AcceptInvitation(ctx, invitation.ID)
	require.NoError(t, err, "should not have failed to accept invite for %v", invitation.ID)

	// should have no pending invitations now
	invitations, _, err = client.Users.ListInvitations(ctx)
	require.NoError(t, err, "could not list invitations in repo %s", fullName)
	require.Empty(t, invitations, "should not have any invitations")
}
