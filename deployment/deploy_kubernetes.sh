#!/bin/bash

#allowed values: mini, stag, prod
SERVER_TYPE=${1:-mini}

echo "Deploying on server: "$SERVER_TYPE

eval $(minikube docker-env)

AWS_ECR="755455355830.dkr.ecr.us-east-2.amazonaws.com"

echo "Removing previous deployment "
echo "---------------------------- "
cd ..
for d in */ ; do
  if [[ $d = *"svc"* ]]; then
    echo "---------------------------- "
    echo "Processing: $d"
    echo "---------------------------- "
    kubectl delete deployment ${d%-svc/}-depl
    kubectl delete service ${d%-svc/}-srvc 
  fi
done

#kubectl delete service database-service
#kubectl delete secret db-secret

#kubectl delete deployment kong
#kubectl delete service kong-admin
#kubectl delete service kong-admin-ssl
#kubectl delete service kong-proxy-ssl
#kubectl delete service kong-proxy

#kubectl create -f db.yaml
#kubectl create -f secret.yaml

echo "Deploying "
echo "--------- "
for d in */ ; do
  if [[ $d = *"svc"* ]]; then
    echo "---------------------------- "
    echo "Processing: $d"
    echo "---------------------------- "
    cd $d
    pwd 
    make build
    make reverse
    make AWS_ECR="$AWS_ECR" REPOSITORY_NAME="${d%/}" docker 
    kubectl create -f ${d%-svc/}-$SERVER_TYPE.yaml
    cd ..
  fi
done

