#!/usr/bin/env sh

env > .env

echo "processing config map..."
oc process -f ./oldMozart/ConfigMap.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing secret..."
oc process -f ./oldMozart/Secret.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing service..."
oc process -f ./oldMozart/Service.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing deployment config..."
oc process -f ./oldMozart/DeploymentConfig.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "processing route..."
oc process -f ./oldMozart/Route.yaml --param-file=.env --ignore-unknown-parameters=true | oc apply -f -

echo "rolling out $SERVICE_NAME ..."
oc rollout latest dc/$SERVICE_NAME
