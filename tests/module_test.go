package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gitCommitHash string

func getGitCommitHash(t *testing.T) string {
	r, err := git.PlainOpen("../")
	require.NoError(t, err)

	h, err := r.ResolveRevision(plumbing.Revision("HEAD"))
	require.NoError(t, err)

	return h.String()[0:7]
}

func TestModule(t *testing.T) {
	gitCommitHash = getGitCommitHash(t)

	initOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
	})
	terraform.Init(t, initOptions)

	t.Run("test-1", ModuleTest1)
}

func ModuleTest1(t *testing.T) {
	t.Parallel()

	options := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		Vars: map[string]interface{}{
			"policy_json": `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Action": [
				"ec2:CreateVpc",
				"ec2:DeleteVpc"
			],
			"Resource": "arn:aws:ec2:*:*:vpc/test-1"
		}
	]
}`,
			"test_name": fmt.Sprintf("terraform-iam-role-module-test-%s-1", gitCommitHash),
		},
	})

	terraform.WorkspaceSelectOrNew(t, options, "test-1")

	_, err := terraform.ApplyAndIdempotentE(t, options)
	if !testing.Short() || err == nil {
		defer terraform.Destroy(t, options)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	require.NoError(t, err)
	cfg.Region = "ca-central-1"

	iamClient := iam.NewFromConfig(cfg)

	roleOutput, err := iamClient.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: aws.String(fmt.Sprintf("terraform-iam-role-module-test-%s-1", gitCommitHash)),
	})
	assert.NoError(t, err)

	assert.NotNil(t, roleOutput.Role.AssumeRolePolicyDocument)
	assert.Contains(t, *roleOutput.Role.AssumeRolePolicyDocument, "sts%3AAssumeRole")
	assert.Contains(t, *roleOutput.Role.AssumeRolePolicyDocument, fmt.Sprintf("aws%%3APrincipalTag%%2Fassume-terraform-iam-role-module-test-%s-1", gitCommitHash))

	rolePolicyOutput, err := iamClient.GetRolePolicy(context.TODO(), &iam.GetRolePolicyInput{
		RoleName:   roleOutput.Role.RoleName,
		PolicyName: roleOutput.Role.RoleName,
	})

	assert.NoError(t, err)

	policyDocument := rolePolicyOutput.PolicyDocument

	assert.NotNil(t, policyDocument)
	assert.Contains(t, *policyDocument, "vpc%2Ftest-1")
}
