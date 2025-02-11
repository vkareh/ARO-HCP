using '../templates/svc-cluster.bicep'

param kubernetesVersion = '{{ .kubernetesVersion }}'
param istioVersions = '{{ .svc.istio.versions }}'
param vnetAddressPrefix = '{{ .vnetAddressPrefix }}'
param subnetPrefix = '{{ .subnetPrefix }}'
param podSubnetPrefix = '{{ .podSubnetPrefix }}'
param aksClusterName = '{{ .aksName }}'
param aksKeyVaultName = '{{ .svc.etcd.kvName }}'
param aksEtcdKVEnableSoftDelete = {{ .svc.etcd.kvSoftDelete }}

param userAgentMinCount = {{ .svc.userAgentPool.minCount }}
param userAgentMaxCount = {{ .svc.userAgentPool.maxCount }}
param userAgentVMSize = '{{ .svc.userAgentPool.vmSize }}'
param aksUserOsDiskSizeGB = {{ .svc.userAgentPool.osDiskSizeGB }}
param userAgentPoolAZCount = {{ .svc.userAgentPool.azCount }}

param disableLocalAuth = {{ .frontend.cosmosDB.disableLocalAuth }}
param deployFrontendCosmos = {{ .frontend.cosmosDB.deploy }}
param rpCosmosDbName = '{{ .frontend.cosmosDB.name }}'
param rpCosmosDbPrivate = {{ .frontend.cosmosDB.private }}

param maestroMIName = '{{ .maestro.server.managedIdentityName }}'
param maestroNamespace = '{{ .maestro.server.k8s.namespace }}'
param maestroServiceAccountName = '{{ .maestro.server.k8s.serviceAccountName }}'
param maestroEventGridNamespacesName = '{{ .maestro.eventGrid.name }}'
param maestroServerMqttClientName = '{{ .maestro.server.mqttClientName }}'
param maestroCertDomain = '{{ .maestro.certDomain}}'
param maestroPostgresServerName = '{{ .maestro.postgres.name }}'
param maestroPostgresServerMinTLSVersion = '{{ .maestro.postgres.minTLSVersion }}'
param maestroPostgresServerVersion = '{{ .maestro.postgres.serverVersion }}'
param maestroPostgresServerStorageSizeGB = {{ .maestro.postgres.serverStorageSizeGB }}
param maestroPostgresDatabaseName = '{{ .maestro.postgres.databaseName }}'
param deployMaestroPostgres = {{ .maestro.postgres.deploy }}
param maestroPostgresPrivate = {{ .maestro.postgres.private }}

param csPostgresDeploy = {{ .clusterService.postgres.deploy }}
param csPostgresServerName = '{{ .clusterService.postgres.name }}'
param csPostgresServerMinTLSVersion = '{{ .clusterService.postgres.minTLSVersion }}'
param clusterServicePostgresPrivate = {{ .clusterService.postgres.private }}

param serviceKeyVaultName = '{{ .serviceKeyVault.name }}'
param serviceKeyVaultResourceGroup = '{{ .serviceKeyVault.rg }}'

param acrPullResourceGroups = ['{{ .global.rg }}']
param clustersServiceAcrResourceGroupNames = ['{{ .clusterService.acrRG }}']
param useCustomACRTokenManagementRole = {{ .global.manageTokenCustomRole }}

param oidcStorageAccountName = '{{ .oidcStorageAccountName }}'
param aroDevopsMsiId = '{{ .aroDevopsMsiId }}'

param regionalCXDNSZoneName = '{{ .dns.regionalSubdomain }}.{{ .dns.cxParentZoneName }}'

param regionalResourceGroup = '{{ .regionRG }}'

param frontendIngressCertName = '{{ .frontend.cert.name }}'

// Azure Monitor Workspace
param azureMonitorWorkspaceName = '{{ .monitoring.workspaceName }}'
