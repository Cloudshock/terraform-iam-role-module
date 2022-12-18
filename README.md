# terraform-iam-role-module

<!-- BEGIN_TF_DOCS -->
This module creates an IAM Role with a standard Trust Relationship policy
that is based on a aws:PrincipalTag condition. It is intended to be
incorporated into a specific Terraform Module that specifies the IAM policy
that is passed into this module.

## Usage

This example shows how a parent module can embed this module in its
configuration.

```hcl
data "aws_iam_policy_document" "this" {
  statement {
    actions = [
      "ec2:CreateVpc",
      "ec2:DeleteVpc",
    ]
    resources = [
      "arn:aws:ec2:*:*:vpc/*"
    ]
  }
}

module "terraform_iam_role" {
  source = "app.terraform.io/cloudshock/terraform-iam-role/aws"
  version = "0.1.0"

  account_id      = var.account_id
  policy_json     = data.aws_iam_policy_document.this.json
  repository_name = "example"
}
```

In the above example, an IAM Role named example will be created with a Trust
Relationship that allows any Principal from the specified account with the
tag `assume-example` set to `true`. The role will also have an inline policy
named **example** created with the provided JSON encoded policy.

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.46.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | >= 4.46.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_iam_role.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_policy_document.assume_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_account_id"></a> [account\_id](#input\_account\_id) | The AWS account ID number that is used in IAM policies. | `string` | n/a | yes |
| <a name="input_policy_json"></a> [policy\_json](#input\_policy\_json) | The JSON encoded inline policy to create with the role. | `string` | n/a | yes |
| <a name="input_repository_name"></a> [repository\_name](#input\_repository\_name) | The name of the repository to which this role is linked. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_role_arn"></a> [role\_arn](#output\_role\_arn) | The ARN of the role created by this module. |
<!-- END_TF_DOCS -->