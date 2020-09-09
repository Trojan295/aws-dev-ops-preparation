# AWS Serverless Application Model Lambda Notes example application

## How to run it

1. Deploy the SAM package
```
cd backend
make deploy
```
2. Setup the repository in the created AWS Amplify app. You can use the AWS Console, so it will guide you through this
3. Browse to the created Amplify app webpage

## Lambda CodeDeploy deployments

SAM uses Lambda aliases to perform canary and linear deployments. This allows to gradually shift traffic to the new version and rollback in case errors are appearing.

This example shows two options how to perform checks on a new Lambda deployment and rollback it in case of failures:

1. Rollback triggered by Cloudwatch Alarms in the `AllowTraffic` stage. The following setup triggers a rollback of the deployment in case there are more that 5% of 5xx errors for 1 minute on the API Gateway during traffic shifting. Note that it means 5% of all Lambda calls and as SAM uses Lambda weighted alias you have to take it into account. For ex. if using `Canary10Percent5Minutes`, then 10% of the traffic is directed to the new Lambda, so with 5% error rate, 50% of the traffic to the new Lambda would need to error to trigger the rollback
```
  NotesApiServerErrorAlarm:
    Type: AWS::CloudWatch::Alarm
    Properties:
      AlarmName: NotesApiServerErrorAlarm
      EvaluationPeriods: 1
      Metrics:
      - Id: m1
        MetricStat:
          Metric:
            Dimensions:
              - Name: ApiName
                Value: !Ref ApiName
            MetricName: 5XXError
            Namespace: AWS/ApiGateway
          Period: !!int 60
          Stat: Average
      ComparisonOperator: GreaterThanThreshold
      Threshold: 0.05
      TreatMissingData: notBreaching

  AddNote:
    Type: AWS::Serverless::Function
    Properties:
      DeploymentPreference:
        Type: Canary10Percent5Minutes
        Alarms:
          - !Ref NotesApiServerErrorAlarm
      [...]
```

2. Hook triggered rollback. You can call a Lambda function in the `BeforeAllowTraffic` or `AfterAllowTraffic` stage, to verify the new Lambda function or run final tests. You muist call the `codedeploy:PutLifecycleEventHookExecutionStatus` API call to tell CodeDeploy about the status of the deployment. CodeDeploy will fail the deployment after 1 hour, if no call is made
```
  GetNotes:
    Type: AWS::Serverless::Function
    Properties:
      DeploymentPreference:
        Hooks:
          PostTraffic: !Ref ValidateAPI
      [...]

  ValidateAPI:
    Type: AWS::Serverless::Function
    Properties:
      FunctionName: CodeDeployHook_ValidateAPI
      Policies:
        - Version: '2012-10-17' 
          Statement:
            - Effect: Allow
              Action:
                - codedeploy:PutLifecycleEventHookExecutionStatus 
              Resource:
                !Sub 'arn:aws:codedeploy:${AWS::Region}:${AWS::AccountId}:deploymentgroup:${ServerlessDeploymentApplication}/*'
      [...]
```
Note that the hook Lambda function needs the IAM permissions to call the AWS API

## Notes about SAM

1. Hook functions must start with prefix "CodeDeployHook_" or you have to provide an custom IAM role for the CodeDeploy in `DeploymentPreference`
2. The hook functions must response to AWS API, if the hook passed or failed. CodeDeploy timeouts after 1 hour waiting resulting in a fail. You need to provide the policy to those functions, so they can all the API
3. The `Alarms` in `DeploymentPreference` can be used to rollback the deployment on Cloudwatch Alarm. The AWS docs suggest the other way - that they are triggered by a failed deployment