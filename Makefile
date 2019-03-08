ENV_FILE := .env
include ${ENV_FILE}
export $(shell sed 's/=.*//' ${ENV_FILE})


oc_login:
	${OC} login ${OC_URL} -u ${OC_USER} -p ${OC_PASSWORD} --insecure-skip-tls-verify=true

create-image:
	cd docker/
	docker build . -t quay.io/redhatdemo/demo4-load-test-websocket:latest

deploy-load-test: oc_login
    oc new-project demo4-load-test
	oc process -f openshift/deployment --param DURATION=${DURATION} --param USERS_PER_POD=${USERS_PER_POD} --param REPLICAS=${REPLICAS} --param SOCKET_ADDRESS=${SOCKET_ADDRESS} | oc create -f - 

remove-load-test: oc_login
	oc delete project demo4-load-test