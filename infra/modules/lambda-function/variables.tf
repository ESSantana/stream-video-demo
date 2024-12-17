variable function_name {    
  type        = string
  description = "Lambda function name"
}

variable stage {    
  type        = string
  default     = "development"
  description = "Stage (environment) of the application"
}

variable aws_region {     
  type        = string
  default     = "sa-east-1"
  description = "AWS region"
}

variable handler {    
  type        = string
  default     = "handler"
  description = "Lambda function handler"
}

variable runtime {    
  type        = string
  default     = "provided.al2023"
  description = "Lambda function runtime"
}

variable architecture {    
  type        = set(string)
  default     = ["arm64"]
  description = "Lambda function architecture"
}

variable role_arn {    
  type        = string
  default     = ""
  description = "IAM role ARN for the lambda function"
}

variable memory_size {    
  type        = number
  default     = 128
  description = "Lambda function memory size in megabytes"
}

variable environment_variables {    
  type        = map(string)
  default     = {}
  description = "Environment variables for the lambda function"
}

variable reserved_concurrent_executions {    
  type        = number
  default     = -1
  description = "The number of simultaneous executions to reserve for the function"
}

variable tags {    
  type        = map(string)
   default    = {}
  description = "Tags for the lambda function"
}