module "new_lambda_test" { 
    source = "./modules/lambda-function"

    function_name = "new-lambda-test"
    stage         = var.stage
}