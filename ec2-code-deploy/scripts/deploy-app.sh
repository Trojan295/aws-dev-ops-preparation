#!/bin/bash

set -e

. scripts/vars.sh

SUFFIX=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 16 | head -n 1)

STACK_OUTPUTS=$(aws cloudformation describe-stacks --stack-name ec2-deploy-stack | jq -r '.Stacks[0].Outputs[]')

APPLICATION_NAME=$(echo "$STACK_OUTPUTS" | jq -r 'select(.OutputKey=="CodeDeployApplicationName") | .OutputValue')
REVISION_BUCKET=$(echo "$STACK_OUTPUTS" | jq -r 'select(.OutputKey=="RevisionBucketName") | .OutputValue')
DEPLOYMENT_GROUP_NAME=$(echo "$STACK_OUTPUTS" | jq -r 'select(.OutputKey=="DeploymentGroupName") | .OutputValue')

REVISION_KEY="revisions/${APPLICATION_NAME}-${SUFFIX}.zip"

aws deploy push \
	--application-name $APPLICATION_NAME \
	--ignore-hidden-files \
	--s3-location s3://${REVISION_BUCKET}/${REVISION_KEY} \
	--source .

aws deploy create-deployment \
	--application-name $APPLICATION_NAME \
	--s3-location bucket=${REVISION_BUCKET},key=${REVISION_KEY},bundleType=zip \
	--deployment-group-name $DEPLOYMENT_GROUP_NAME
