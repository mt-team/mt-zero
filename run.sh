#!/bin/bash

app=$1
env=$2

if [[  $app == "default" ]]; then
  echo "Wrong app name."
  exit 1
fi

#./bin/$app  -f ./src/${app}/etc/${app}-local.yaml
if [[  $env == "test" ]]; then
  ./$app  -f ./${app}-test.yaml
elif [[  $env == "release" ]]; then
  ./$app  -f ./${app}.yaml
else
  ./bin/$app  -f ./src/${app}/etc/${app}-local.yaml
fi
