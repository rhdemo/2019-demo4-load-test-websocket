ENV_FILE := .env
include ${ENV_FILE}
export $(shell sed 's/=.*//' ${ENV_FILE})


oc_login:
	${OC} login ${OC_URL} -u ${OC_USER} -p ${OC_PASSWORD} --insecure-skip-tls-verify=true

create-image:
	@echo Create Image
	cd docker/ && docker build . -t ${IMAGE}

push-image:
	@echo Push Image
	docker push ${IMAGE}


deploy-load-tester: oc_login
	@echo Deploying Jobs for Load Test	
	ansible-playbook openshift/playbook.yml -e users=${USERS} -e namespace=${NAMESPACE} -e duration=${DURATION} -e image=${IMAGE} -e replicas=${REPLICAS} -e ws_address=${SOCKET_ADDRESS} -e movement=${MOVEMENT} -v

clean-namespace: oc_login
	@echo Removing Load Test Assets
	oc delete pods all -n ${NAMESPACE} --ignore-not-found=true
	oc delete configmap load-configmap -n ${NAMESPACE} --ignore-not-found=true

remove-load-test: oc_login
	@echo Removing Load Test Namespace
	oc delete project ${NAMESPACE} --ignore-not-found=true
