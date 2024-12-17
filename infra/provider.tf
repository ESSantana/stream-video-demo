provider "aws" { 
    region = var.aws_region
} 

provider archive  {
    source = "hashicorp/archive"
}