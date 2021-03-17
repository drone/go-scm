package fake_test

import (
	"context"
	"github.com/jenkins-x/go-scm/scm"
	"testing"

	"github.com/jenkins-x/go-scm/scm/driver/fake"
	"github.com/stretchr/testify/require"
)

func TestDeploy(t *testing.T) {
	client, _ := fake.NewDefault()

	ctx := context.Background()

	repo := "myorg/myrepo"

	AssertDeploymentSize(t, ctx, client, 0, repo)

	// lets create a deployment
	input := &scm.DeploymentInput{
		Ref:         "12345",
		Payload:     "deplyment",
		Environment: "production",
		Description: "my funky deployment thingy",
	}
	deploy, _, err := client.Deployments.Create(ctx, repo, input)
	require.NoError(t, err, "failed to create deploy in repo %s", repo)
	require.NotNil(t, deploy, "should have created a deployment")

	AssertDeploymentSize(t, ctx, client, 1, repo)
	AssertDeploymentStatusSize(t, ctx, client, 0, repo, deploy.ID)

	deploy2, _, err := client.Deployments.Find(ctx, repo, deploy.ID)
	require.NoError(t, err, "failed to find deploy in repo %s for deployment %s ", repo, deploy.ID)
	require.NotNil(t, deploy2, "should have found a deployment in repo %s for deployment %s ", repo, deploy.ID)

	// lets create a status
	statusInput := &scm.DeploymentStatusInput{
		State:       "success",
		TargetLink:  "https://acme.com/my/app/thing",
		Description: "my update",
		Environment: "production",
	}
	status, _, err := client.Deployments.CreateStatus(ctx, repo, deploy.ID, statusInput)
	require.NoError(t, err, "failed to create status in repo %s for status %s", repo, deploy.ID)
	require.NotNil(t, status, "should have created a status")

	AssertDeploymentStatusSize(t, ctx, client, 1, repo, deploy.ID)

	status2, _, err := client.Deployments.FindStatus(ctx, repo, deploy.ID, status.ID)
	require.NoError(t, err, "failed to find status in repo %s for deployment %s status %s", repo, deploy.ID, status.ID)
	require.NotNil(t, status2, "should have found a status in repo %s for deployment %s status %s", repo, deploy.ID, status.ID)

	// lets delete the deployment
	_, err = client.Deployments.Delete(ctx, repo, deploy.ID)
	require.NoError(t, err, "failed to delete deploy in repo %s", repo)

	AssertDeploymentSize(t, ctx, client, 0, repo)
}

func AssertDeploymentSize(t *testing.T, ctx context.Context, client *scm.Client, size int, repo string) []*scm.Deployment {
	deploys, _, err := client.Deployments.List(ctx, repo, scm.ListOptions{})
	require.NoError(t, err, "could not list deploys in repo %s", repo)
	require.Len(t, deploys, size, "deploy size")
	return deploys
}

func AssertDeploymentStatusSize(t *testing.T, ctx context.Context, client *scm.Client, size int, repo, deploymentID string) []*scm.DeploymentStatus {
	statuses, _, err := client.Deployments.ListStatus(ctx, repo, deploymentID, scm.ListOptions{})
	require.NoError(t, err, "could not list statuses in repo %s deploymentID %s", repo, deploymentID)
	require.Len(t, statuses, size, "status size")
	return statuses
}
