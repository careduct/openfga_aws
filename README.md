# OpenFGA HTTP Endpoint Deployment on AWS Lambda and Amazon Aurora Postgres serverless v2

Welcome to our repository! Here, you'll find all the necessary code to deploy the HTTP endpoint of OpenFGA on AWS Lambda and Amazon RDS. This guide will walk you through the steps to get everything up and running.

## Overview
OpenFGA is a versatile and open-source service designed to help developers with access control management. By integrating OpenFGA's HTTP endpoint with AWS API Gateway and Lambda, you can achieve scalable and serverless function execution. Additionally, using Amazon RDS for database management ensures a strong and scalable foundation for handling authorization tasks.

The process involves creating a Lambda extension that initializes OpenFGA during the cold start, setting up both gRPC and HTTP servers. For every incoming request, the AWS API Gateway directs the request to the server operating on the Lambda extension. This setup facilitates efficient and scalable authorization services.

## Prerequisites

Before you begin, make sure you have the following:
- An AWS account
- AWS CLI installed and configured
- Basic knowledge of AWS services (Lambda, API Gateway, RDS, IAM)
- Familiarity with OpenFGA

## Deployment Steps

1. **Prepare Your AWS Environment:**
   - create a provider.tf file at terraform folder and provide the credentials and regions where you want the infra deployed
   - run terraform init

2. **Deploy the Lambda Function:**
   - run 'make package' on the services/openfga_http folder
   - run terraform apply

3. **Test Your Deployment:**
   - As output you will get an url of the API Gateway endpoint, ex: https://z6u4i2vz84.execute-api.eu-central-1.amazonaws.com/openfga
   - Add 'stores' on the path to get the usable Openfga http API, ex: https://z6u4i2vz84.execute-api.eu-central-1.amazonaws.com/openfga/stores 

## Additional Resources

- [AWS Lambda Documentation](https://docs.aws.amazon.com/lambda/)
- [Amazon RDS Documentation](https://docs.aws.amazon.com/rds/)
- [OpenFGA Documentation](https://docs.openfga.org/)

## Contributing

We welcome contributions! Please feel free to submit pull requests or open issues to suggest improvements or add new features.

---

Thank you for visiting our repository. We hope this guide helps you deploy OpenFGA's HTTP endpoint seamlessly on AWS Lambda and RDS.
