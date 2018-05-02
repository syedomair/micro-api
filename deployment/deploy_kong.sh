#!/bin/bash

eval $(minikube docker-env)

KONG_ADMIN="$(minikube service kong-admin --url)"
PUBLIC_SRVC="$(minikube service public-srvc --url)"
ROLES_SRVC="$(minikube service roles-srvc --url)"
USERS_SRVC="$(minikube service users-srvc --url)"
BATCH_SRVC="$(minikube service batch-tasks-srvc --url)"
echo kong-admin:$KONG_ADMIN
echo public-srvc:$PUBLIC_SRVC
echo roles-srvc:$ROLES_SRVC
echo users-srvc:$USERS_SRVC

LENGTH="$(curl -X GET --url $KONG_ADMIN/routes | jq ' .["data"] | length') "
echo "Length:"$LENGTH

ROUTES="$(curl -X GET --url $KONG_ADMIN/routes )"
for (( c=0; c<=$LENGTH-1; c++ ))
do  
   echo "--------------------"
   echo "Deleting route $c "
   echo "--------------------"
   ROUTE_ID="$(jq -n "$ROUTES" | jq -c '.[]['$c'].id' | awk -F '"' '{print $2}' | tr -dc '[[:print:]]' )"
   URL=$KONG_ADMIN/routes/$ROUTE_ID 
   echo $URL
   curl -i -X DELETE --url $URL
done

LENGTH="$(curl -X GET --url $KONG_ADMIN/services | jq ' .["data"] | length') "
echo "Length:"$LENGTH

SERVICES="$(curl -X GET --url $KONG_ADMIN/services )"
for (( c=0; c<=$LENGTH-1; c++ ))
do  
   echo "--------------------"
   echo "Deleting services $c "
   echo "--------------------"
   SERVICE_ID="$(jq -n "$SERVICES" | jq -c '.[]['$c'].id' | awk -F '"' '{print $2}' | tr -dc '[[:print:]]' )"
   URL=$KONG_ADMIN/services/$SERVICE_ID 
   echo $URL
   curl -i -X DELETE --url $URL
done

deploySecure () {
    echo "-----------------------------"
    echo Deploying $1
    echo
    echo "-----------------------------Creating new Service for $1 -----------------------------"
    curl -i -X POST --url $KONG_ADMIN/services/ --data 'name='$1 --data 'url='$3$2
    echo
    echo "-----------------------------Creating new Route for $1 -----------------------------"
    curl -i -X POST --url $KONG_ADMIN/services/$1/routes --data 'paths[]='$2 --data 'methods[]=POST' --data 'methods[]=GET' --data 'methods[]=PATCH' --data 'methods[]=DELETE'
    echo
    echo "-----------------------------Applying key-auth to service for $1 ----------------------"
    curl -i -X POST --url $KONG_ADMIN/services/$1/plugins --data "name=key-auth"
    echo
}
deployPublic () {
    echo "-----------------------------"
    echo Deploying $1
    echo "-----------------------------"
    echo
    echo "-----------------------------Creating new Service for $1 -----------------------------"
    curl -i -X POST --url $KONG_ADMIN/services/ --data 'name='$1 --data 'url='$3$2
    echo
    echo "-----------------------------Creating new Route for $1 -----------------------------"
    curl -i -X POST --url $KONG_ADMIN/services/$1/routes --data 'paths[]='$2 --data 'methods[]='$4
}

deployPublic "public-register" "/v1/register" $PUBLIC_SRVC "POST"
deployPublic "public-authenticate" "/v1/authenticate" $PUBLIC_SRVC "POST"
deploySecure "users" "/v1/users" $USERS_SRVC 
deploySecure "roles" "/v1/roles" $ROLES_SRVC 
deploySecure "batch-tasks" "/v1/batch/users" $BATCH_SRVC 
