# CodeDeploy to EC2 instances

## How to run it

1. Deploy the CloudFormation stack
```
make cf-deploy
```
2. Deploy the application
```
make app-deploy
```

## SSM State Manager

This example uses SSM State Manager to configure the EC2 instances, instead of userdata or cf-init scripts. It creates two associations:
- to install CodeDeploy agent
- to install httpd server

The only requirements to use SSM is to install the SSM agent on the instance (Amazon Linux 2 AMIs have it already installed) and to attach an IAM role to the instance, which can make API calls to SSM.

## CodeDeploy

In this example a basic AppSpec file is used. It just deploys the `index.html` file to https document root and does not define any lifecycle hooks. The deployment is automated in [deploy-app.sh](scripts/deploy-app.sh). It creates the revision and triggers a CodeDeploy deployment of the application.
