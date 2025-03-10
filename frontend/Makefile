-include ../setup-env.mk
-include ../helm-cmd.mk
HELM_CMD ?= helm upgrade --install

ifndef COMMIT
COMMIT := $(shell git rev-parse --short=7 HEAD)
endif
ARO_HCP_IMAGE_REGISTRY ?= ${ARO_HCP_IMAGE_ACR}.azurecr.io
ARO_HCP_IMAGE_REPOSITORY ?= arohcpfrontend
ARO_HCP_FRONTEND_IMAGE ?= $(ARO_HCP_IMAGE_REGISTRY)/$(ARO_HCP_IMAGE_REPOSITORY):$(COMMIT)

.DEFAULT_GOAL := frontend

frontend:
	go build -o aro-hcp-frontend .

info:
	@echo "ARO_HCP_FRONTEND_IMAGE: ${ARO_HCP_FRONTEND_IMAGE}"
	@echo "COMMIT: ${COMMIT}"

run:
	./aro-hcp-frontend --use-cache --location ${LOCATION} \
		--clusters-service-url http://localhost:8000 \
		--cluster-service-provision-shard 1 \
		--cluster-service-noop-provision \
		--cluster-service-noop-deprovision

clean:
	rm -f aro-hcp-frontend
.PHONY: clean

build-push: image push

image:
	pushd .. && git archive --output frontend/archive.tar.gz HEAD && popd
	docker build -f "./Dockerfile" -t ${ARO_HCP_FRONTEND_IMAGE} .
	rm -f archive.tar.gz
.PHONY: image

push: image
	az acr login --name ${ARO_HCP_IMAGE_ACR}
	docker push ${ARO_HCP_FRONTEND_IMAGE}

deploy:
	FRONTEND_MI_CLIENT_ID=$$(az identity show \
			-g ${RESOURCEGROUP} \
			-n frontend \
			--query clientId -o tsv) && \
	FRONTEND_MI_TENANT_ID=$$(az identity show \
			-g ${RESOURCEGROUP} \
			-n frontend \
			--query tenantId -o tsv) &&\
	IMAGE_PULLER_MI_CLIENT_ID=$$(az identity show \
			-g ${RESOURCEGROUP} \
			-n image-puller \
			--query clientId -o tsv) && \
	IMAGE_PULLER_MI_TENANT_ID=$$(az identity show \
			-g ${RESOURCEGROUP} \
			-n image-puller \
			--query tenantId -o tsv) && \
	SECRET_STORE_MI_CLIENT_ID=$$(az aks show --resource-group ${RESOURCEGROUP} \
			--name ${AKS_NAME} \
			--query addonProfiles.azureKeyvaultSecretsProvider.identity.clientId \
			--output tsv) && \
	DB_URL=$$(az cosmosdb show -n ${DB_NAME} -g ${RESOURCEGROUP} --query documentEndpoint -o tsv) && \
	TENANT_ID=$(shell az account show --query tenantId --output tsv) && \
	kubectl create namespace aro-hcp --dry-run=client -o json | kubectl apply -f - && \
	kubectl label namespace aro-hcp "istio.io/rev=${ISTO_TAG}" --overwrite=true && \
	${HELM_CMD} aro-hcp-frontend-dev \
		deploy/helm/frontend/ \
		--set azure.clientId=$${SECRET_STORE_MI_CLIENT_ID} \
		--set azure.tenantId=$${TENANT_ID} \
		--set configMap.databaseName=${DB_NAME} \
		--set configMap.databaseUrl="$${DB_URL}" \
		--set configMap.frontendMiClientId="$${FRONTEND_MI_CLIENT_ID}" \
		--set credsKeyVault.name=${SERVICE_KEY_VAULT} \
		--set credsKeyVault.secret=${CERTIFICATE_NAME} \
		--set serviceAccount.workloadIdentityClientId="$${FRONTEND_MI_CLIENT_ID}" \
		--set serviceAccount.workloadIdentityTenantId="$${FRONTEND_MI_TENANT_ID}" \
		--set pullBinding.workloadIdentityClientId="$${IMAGE_PULLER_MI_CLIENT_ID}" \
		--set pullBinding.workloadIdentityTenantId="$${IMAGE_PULLER_MI_TENANT_ID}" \
		--set configMap.currentVersion=${ARO_HCP_FRONTEND_IMAGE} \
		--set configMap.location=${LOCATION}  \
		--set deployment.imageName=${ARO_HCP_FRONTEND_IMAGE} \
		--set pullBinding.registry=${ARO_HCP_IMAGE_REGISTRY} \
		--set pullBinding.scope=repository:${ARO_HCP_IMAGE_REPOSITORY}:pull \
		--namespace aro-hcp
.PHONY: deploy

undeploy:
	helm uninstall aro-hcp-frontend-dev --namespace aro-hcp
.PHONY: undeploy

smoke-tests:
	go test -v -count 1 ./utils/frontend_smoke_test.go
