terraform {
  required_version = ">= 1.3.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.46.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}

data "aws_caller_identity" "current" {}

module "mut" {
  source = "../"

  account_id      = data.aws_caller_identity.current.account_id
  policy_json     = var.policy_json
  repository_name = var.test_name
}