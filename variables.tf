variable "account_id" {
  type        = string
  description = "The AWS account ID number that is used in IAM policies."
}

variable "policy_json" {
  type        = string
  description = "The JSON encoded inline policy to create with the role."
}

variable "repository_name" {
  type        = string
  description = "The name of the repository to which this role is linked."
}
