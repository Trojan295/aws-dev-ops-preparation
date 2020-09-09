#!/bin/bash

set -e

. ./scripts/vars.sh

CF_ARGS="""--stack-name ${STACK_NAME}
    --template-file ${TEMPLATE_FILE}
	--capabilities CAPABILITY_IAM 
	--parameter-overrides VpcId=${VPC_ID} Subnets=${SUBNET_IDS}"""

aws cloudformation deploy ${CF_ARGS}
