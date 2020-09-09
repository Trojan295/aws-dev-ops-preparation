STACK_NAME=ec2-deploy-stack

VPC_ID=vpc-0835d871
SUBNET_IDS="subnet-06530d60\,subnet-4389a70b"

TEMPLATE_FILE="file://template.yml"

.PHONY: deploy destroy update

deploy:
	aws cloudformation create-stack \
		--stack-name ${STACK_NAME} \
		--template-body ${TEMPLATE_FILE} \
		--capabilities CAPABILITY_IAM \
		--parameters ParameterKey=VpcId,ParameterValue=${VPC_ID} \
		ParameterKey=Subnets,ParameterValue=${SUBNET_IDS}

update:
	aws cloudformation update-stack \
		--stack-name ${STACK_NAME} \
		--template-body ${TEMPLATE_FILE} \
		--capabilities CAPABILITY_IAM \
		--parameters ParameterKey=VpcId,ParameterValue=${VPC_ID} \
		ParameterKey=Subnets,ParameterValue=${SUBNET_IDS}

destroy:
	aws cloudformation delete-stack --stack-name ${STACK_NAME}

push_revision:
	aws deploy push \
		--application-name ec2-webpage \
		--ignore-hidden-files \
		--s3-location s3://dczaja-ec2-webpage-bucket/revisions/ec2-webpage.zip \
		--source .
