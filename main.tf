/*
 * This module creates an IAM Role with a standard Trust Relationship policy
 * that is based on a aws:PrincipalTag condition. It is intended to be
 * incorporated into a specific Terraform Module that specifies the IAM policy
 * that is passed into this module.
 *
 * ## Usage
 *
 * This example shows how a parent module can embed this module in its
 * configuration.
 *
 * ```hcl
 * data "aws_iam_policy_document" "this" {
 *   statement {
 *     actions = [
 *       "ec2:CreateVpc",
 *       "ec2:DeleteVpc",
 *     ]
 *     resources = [
 *       "arn:aws:ec2:*:*:vpc/*"
 *     ]
 *   }
 * }
 *
 * module "terraform_iam_role" {
 *   source = "app.terraform.io/cloudshock/terraform-iam-role/aws"
 *   version = "0.1.0"
 *
 *   account_id      = var.account_id
 *   policy_json     = data.aws_iam_policy_document.this.json
 *   repository_name = "example"
 * }
 * ```
 *
 * In the above example, an IAM Role named example will be created with a Trust
 * Relationship that allows any Principal from the specified account with the
 * tag `assume-example` set to `true`. The role will also have an inline policy
 * named **example** created with the provided JSON encoded policy.
 */

terraform {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.46.0"
    }
  }
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = [
      "sts:AssumeRole",
    ]
    effect = "Allow"
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::${var.account_id}:root",
      ]
    }
    condition {
      test     = "StringEquals"
      variable = "aws:PrincipalTag/assume-${var.repository_name}"
      values = [
        "true",
      ]
    }
  }
}

resource "aws_iam_role" "this" {
  name = var.repository_name

  assume_role_policy = data.aws_iam_policy_document.assume_role.json
  inline_policy {
    name   = var.repository_name
    policy = var.policy_json
  }
}
