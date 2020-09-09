#!/bin/bash

set -e

. scripts/vars.sh

aws cloudformation delete-stack --stack-name $STACK_NAME
