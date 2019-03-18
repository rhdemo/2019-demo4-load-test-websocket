ENV_FILE := .env
include ${ENV_FILE}
export $(shell sed 's/=.*//' ${ENV_FILE})


oc_login:
	${OC} login ${OC_URL} -u ${OC_USER} -p ${OC_PASSWORD} --insecure-skip-tls-verify=true

create-image:
	@echo Create Image
	cd docker/ && docker build . -t quay.io/redhatdemo/demo4-load-test-websocket:latest

push-image:
	@echo Push Image
	docker push quay.io/redhatdemo/demo4-load-test-websocket:latest


deploy-load-test: oc_login
	@echo Deploying Load Test
	oc process -f openshift/deployment.yml --param DURATION=${DURATION} --param USERS=${USERS} --param SOCKET_ADDRESS=${SOCKET_ADDRESS} | oc create -f - 

remove-load-test: oc_login
	@echo Removing Load Test
	oc delete project demo4-load-test --ignore-not-found=true
