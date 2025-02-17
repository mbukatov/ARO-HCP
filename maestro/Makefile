SHELL = /bin/bash

RESOURCEGROUP ?= aro-hcp-${AKSCONFIG}-$(USER)
MAESTRO_INFRA_RESOURCEGROUP ?= aro-hcp-svc-cluster-$(USER)
CONSUMER_NAME ?= aro-hcp-mgmt-cluster-$(USER)
REGION ?= $(shell az group show -n ${RESOURCEGROUP} --query location -o tsv)
EVENTGRID_ID = $(shell az resource list -g ${MAESTRO_INFRA_RESOURCEGROUP} --resource-type "Microsoft.EventGrid/namespaces" --query "[].id" -o tsv)


deploy-server:
	kubectl apply -f deploy/istio.yml
	MAESTRO_MI_CLIENT_ID=$(shell az identity show \
			-g ${RESOURCEGROUP} \
			-n maestro-server \
			--query clientId) && \
	TENANT_ID=$(shell az account show --query tenantId --output tsv) && \
	MAESTRO_KV_NAME=$(shell az keyvault list --query "[?starts_with(name, 'maestro-kv')].name" -g ${MAESTRO_INFRA_RESOURCEGROUP} --output tsv) && \
	EVENTGRID_CA_BASE64='$(shell cat ../dev-infrastructure/scripts/digicert-global-root-g3.crt | base64 | tr -d '\n')' && \
	EVENTGRID_HOSTNAME=$(shell az resource show --ids ${EVENTGRID_ID} --query properties.topicSpacesConfiguration.hostname -o tsv) && \
	oc process --local -f deploy/openshift-templates/namespace-template.yml \
	  -p ISTIO_VERSION=asm-1-20 | oc apply -f - && \
	oc process --local -f deploy/openshift-templates/db-template.yml | oc apply -f - && \
	oc process --local -f deploy/openshift-templates/service-template.yml \
	  -p IMAGE_REGISTRY=quay.io \
	  -p IMAGE_REPOSITORY=redhat-user-workloads/maestro-rhtap-tenant/maestro/maestro \
	  -p IMAGE_TAG=6583de322fb4518a32f6ab21c0c1f6a144ea65cd \
	  -p DB_SSLMODE=disable \
	  -p DB_SECRET_NAME=maestro-db \
	  -p MAESTRO_MI_CLIENT_ID=$${MAESTRO_MI_CLIENT_ID} \
	  -p MAESTRO_KV_NAME=$${MAESTRO_KV_NAME} \
	  -p EVENTGRID_CA_BASE64=$${EVENTGRID_CA_BASE64} \
	  -p EVENTGRID_HOSTNAME=$${EVENTGRID_HOSTNAME} \
	  -p AZURE_TENANT_ID=$${TENANT_ID} | oc apply -f -

deploy-agent:
	MAESTRO_MI_CLIENT_ID=$(shell az identity show \
			-g ${RESOURCEGROUP} \
			-n maestro-consumer \
			--query clientId) && \
	TENANT_ID=$(shell az account show --query tenantId --output tsv) && \
	MAESTRO_KV_NAME=$(shell az keyvault list --query "[?starts_with(name, 'maestro-kv')].name" -g ${MAESTRO_INFRA_RESOURCEGROUP} --output tsv) && \
	EVENTGRID_CA_BASE64='$(shell cat ../dev-infrastructure/scripts/digicert-global-root-g3.crt | base64 | tr -d '\n')' && \
	EVENTGRID_HOSTNAME=$(shell az resource show --ids ${EVENTGRID_ID} --query properties.topicSpacesConfiguration.hostname -o tsv) && \
	oc process --local -f deploy/openshift-templates/agent-template.yml \
	  -p IMAGE_REGISTRY=quay.io \
	  -p IMAGE_REPOSITORY=redhat-user-workloads/maestro-rhtap-tenant/maestro/maestro \
	  -p IMAGE_TAG=ae149df618cb0812d2072b20658a9cff84c087eb \
	  -p MAESTRO_MI_CLIENT_ID=$${MAESTRO_MI_CLIENT_ID} \
	  -p MAESTRO_KV_NAME=$${MAESTRO_KV_NAME} \
	  -p MAESTRO_AGENT_CONSUMER_NAME=${RESOURCEGROUP} \
	  -p EVENTGRID_CA_BASE64=$${EVENTGRID_CA_BASE64} \
	  -p EVENTGRID_HOSTNAME=$${EVENTGRID_HOSTNAME} \
	  -p AZURE_TENANT_ID=$${TENANT_ID} | oc apply -f -

register-agent:
	./register-maestro-consumer.sh ${CONSUMER_NAME}

# still WIP
register-agent-job:
	oc process --local -f deploy/openshift-templates/agent-registration.yml \
	  -p MAESTRO_AGENT_CONSUMER_NAME=${RESOURCEGROUP} | oc apply -f -

.PHONY: deploy-server
