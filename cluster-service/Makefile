SHELL = /bin/bash

AKSCONFIG ?= svc-cluster
RESOURCEGROUP ?= aro-hcp-${AKSCONFIG}-$(USER)
CONSUMER_NAME ?= aro-hcp-mgmt-cluster-$(USER)
REGION ?= $(shell az group show -n ${RESOURCEGROUP} --query location -o tsv)

deploy:
	DOMAIN=$(shell az network dns zone list -g ${RESOURCEGROUP} --query "[?zoneType=='Public'].name" -o tsv) && \
	sed -e "s/DOMAIN/$${DOMAIN}/g" -e "s/REGION/${REGION}/g" -e "s/CONSUMER_NAME/${CONSUMER_NAME}/g" deploy/mvp-provisioning-shards.yml > deploy/tmp-provisioning-shard.yml
	kubectl apply -f deploy/istio.yml
	oc process --local -f deploy/openshift-templates/arohcp-namespace-template.yml \
	  -p ISTIO_VERSION=asm-1-20 | oc apply -f -
	oc process --local -f deploy/openshift-templates/arohcp-db-template.yml | oc apply -f -
	oc process --local -f deploy/openshift-templates/arohcp-secrets-template.yml \
	  -p PROVISION_SHARDS_CONFIG="$$( base64 -i deploy/tmp-provisioning-shard.yml)" | oc apply -f -
	AZURE_CS_MI_CLIENT_ID=$(shell az identity show \
			-g ${RESOURCEGROUP} \
			-n clusters-service \
			--query clientId) && \
	CS_SERVICE_PRINCIPAL_CREDS_BASE64='$(shell az keyvault secret show --vault-name "service-kv-aro-hcp-dev" --name "aro-hcp-dev-sp-cs" | jq .value -r | base64 | tr -d '\n')' && \
	oc process --local -f deploy/openshift-templates/arohcp-service-template.yml \
	  -p AZURE_CS_MI_CLIENT_ID=$${AZURE_CS_MI_CLIENT_ID} \
	  -p REGION=${REGION} \
	  -p CS_SERVICE_PRINCIPAL_CREDS_BASE64=$${CS_SERVICE_PRINCIPAL_CREDS_BASE64} \
	  -p IMAGE_REGISTRY=devarohcp.azurecr.io \
	  -p IMAGE_REPOSITORY=app-sre/uhc-clusters-service \
	  -p IMAGE_TAG=b1268cf | oc apply -f -

provision-shard:
	@DOMAIN=$(shell az network dns zone list -g ${RESOURCEGROUP} --query "[?zoneType=='Public'].name" -o tsv) && \
	sed -e "s/DOMAIN/$${DOMAIN}/g" -e "s/REGION/${REGION}/g" -e "s/CONSUMER_NAME/${CONSUMER_NAME}/g" deploy/dev-provisioning-shards.yml

.PHONY: deploy
