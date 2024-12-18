variable aws_region {    
  type        = string
  default     = "sa-east-1"
  description = "resource aws region"
}

variable stage {    
  type        = string
  description = "Stage (environment) of the workspace"
}

variable accountId {    
  type        = string
  description = "Account ID of the workspace"
}
