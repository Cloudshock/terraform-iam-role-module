package tests

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestModule(t *testing.T) {
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
			"Resource": "arn:aws:ec2:*:*:vpc/test-1-*"
		}
	]
}`,
			"test_name": "test-1",
		},
	})

	terraform.WorkspaceSelectOrNew(t, options, "test-1")

	_, err := terraform.ApplyAndIdempotentE(t, options)
	if !testing.Short() || err == nil {
		defer terraform.Destroy(t, options)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	assert.NoError(t, err)

	iamClient := iam.NewFromConfig(cfg)

	roleOutput, err := iamClient.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: aws.String("test-1"),
	})
	assert.NoError(t, err)

	assert.NotNil(t, roleOutput.Role.AssumeRolePolicyDocument)
	assert.Contains(t, *roleOutput.Role.AssumeRolePolicyDocument, "sts%3AAssumeRole")
	assert.Contains(t, *roleOutput.Role.AssumeRolePolicyDocument, "aws%3APrincipalTag%2Fassume-test-1")

	rolePolicyOutput, err := iamClient.GetRolePolicy(context.TODO(), &iam.GetRolePolicyInput{
		RoleName:   roleOutput.Role.RoleName,
		PolicyName: roleOutput.Role.RoleName,
	})

	assert.NoError(t, err)

	policyDocument := rolePolicyOutput.PolicyDocument

	assert.NotNil(t, policyDocument)
	assert.Contains(t, *policyDocument, "vpc%2Ftest-1")
}
