STACK_NAME=ec2-deploy-stack

.PHONY: cf-deploy cf-destroy app-deploy

cf-deploy:
	./scripts/deploy-stack.sh

cf-destroy:
	aws cloudformation delete-stack --stack-name ${STACK_NAME}

app-deploy:
	./scripts/deploy-app.sh
