terraform {
  backend "s3" {}

  required_providers {
    aws = {
      source = "hashicorp/aws"
      region = var.aws_region
    }

    archive = {
      source = "hashicorp/archive"
      version = "2.7.0"
    }
     
  }
  
}