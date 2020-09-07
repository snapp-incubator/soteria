#!/usr/bin/env sh

env > .env

echo "processing config map..."
oc process -f ./mozart/ConfigMap.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing secret..."
oc process -f ./mozart/Secret.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing service..."
oc process -f ./mozart/Service.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing deployment config..."
oc process -f ./mozart/DeploymentConfig.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing route..."
oc process -f ./mozart/Route.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "rolling out $SERVICE_NAME ..."
oc rollout latest dc/$SERVICE_NAME
