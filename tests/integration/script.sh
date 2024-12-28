#!/bin/bash
set -euo pipefail
set +x

# This test requires users to set up an AWS access key that has access
# to an Secrets Manager.
# After setting AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY environment
# variables, users also needs to set up a environment variable
# PORTER_TEST_SECRETS_MANAGER_REGION that contains the region of the AWS Secrets Manager.
# Then they can run script like so:
# ./tests/integration/script.sh
# This script assumes users are running it from the root directory of the aws-plugin
# repo

TMP=$(mktemp -d -t tmp.XXXXXXXXXX)
PORTER_HOME=$TMP

cleanup(){
    ret=$?
    echo "EXIT STATUS: $ret"
    rm -rf "$TMP"
    echo "cleaned up test successfully"
    exit "$ret"
}
trap cleanup EXIT

if ! command -v jq 2>&1 /dev/null; then
	echo "jq is required."
	exit 1
fi

authSetup=0
if [ -z ${AWS_ACCESS_KEY_ID} ]; then
    echo "AWS_ACCESS_KEY_ID is required for authentication."
	authSetup=1
fi

if [ -z ${AWS_SECRET_ACCESS_KEY} ]; then
    echo "AWS_SECRET_ACCESS_KEY is required for authentication."
	authSetup=1
fi

if [ $authSetup -eq 1 ]; then
	exit 1
fi

if [ -z $PORTER_TEST_SECRETS_MANAGER_REGION ]; then
    echo "PORTER_TEST_SECRETS_MANAGER_REGION is required for running this test."
	exit 1
fi

export PORTER_HOME=$PORTER_HOME
export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY

mkdir -p $PORTER_HOME/plugins/aws
cp ./bin/plugins/aws/aws $PORTER_HOME/plugins/aws/aws

curl -L https://cdn.porter.sh/latest/install-linux.sh | PORTER_HOME=$PORTER_HOME bash
PORTER_CMD="${PORTER_HOME}/porter --verbosity=debug"
secret_value=super-secret

cp ./tests/integration/testdata/config-test.yaml ${PORTER_HOME}/config.yaml

${PORTER_CMD} plugins list
cd ./tests/integration/testdata && ${PORTER_CMD} install --force --param password=$secret_value

id=$(${PORTER_CMD} installation runs list aws-plugin-test -o json | jq -r '.[].id' | head -1)

if [ -z ${id} ]; then
	echo "failed to get run id"
	exit 1
fi

value=$(aws secretsmanager get-secret-value --region $PORTER_TEST_SECRETS_MANAGER_REGION --secret-id $id-password | jq -r '.SecretString')

if [[ $value == $secret_value ]]
then
	echo "test run successfully"
	exit 0
else
	echo "test failed"
	echo "expected to retrieve value: $secret_value from AWS Secrets Manager, but got: $value"
	exit 1
fi
