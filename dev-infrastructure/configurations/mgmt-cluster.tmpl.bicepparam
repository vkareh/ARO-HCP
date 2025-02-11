using '../templates/mgmt-cluster.bicep'

// AKS
param kubernetesVersion = '{{ .kubernetesVersion}}'
param vnetAddressPrefix = '{{ .vnetAddressPrefix }}'
param subnetPrefix = '{{ .subnetPrefix }}'
param podSubnetPrefix = '{{ .podSubnetPrefix }}'
param aksClusterName = '{{ .aksName }}'
param aksKeyVaultName = '{{ .mgmt.etcd.kvName }}'
param aksEtcdKVEnableSoftDelete = {{ .mgmt.etcd.kvSoftDelete }}
param systemAgentMinCount = {{ .mgmt.systemAgentPool.minCount}}
param systemAgentMaxCount = {{ .mgmt.systemAgentPool.maxCount }}
param systemAgentVMSize = '{{ .mgmt.systemAgentPool.vmSize }}'
param aksSystemOsDiskSizeGB = {{ .mgmt.systemAgentPool.osDiskSizeGB }}
param userAgentMinCount = {{ .mgmt.userAgentPool.minCount }}
param userAgentMaxCount = {{ .mgmt.userAgentPool.maxCount }}
param userAgentVMSize = '{{ .mgmt.userAgentPool.vmSize }}'
param aksUserOsDiskSizeGB = {{ .mgmt.userAgentPool.osDiskSizeGB }}
param userAgentPoolAZCount = {{ .mgmt.userAgentPool.azCount }}

// Maestro
param maestroConsumerName = '{{ .maestro.consumerName }}'
param maestroEventGridNamespacesName = '{{ .maestro.eventGrid.name }}'
param maestroCertDomain = '{{ .maestro.certDomain }}'

// ACR
param acrPullResourceGroups = ['{{ .global.rg }}']

// Region
param regionalResourceGroup = '{{ .regionRG }}'

// CX KV
param cxKeyVaultName = '{{ .cxKeyVault.name }}'

// MSI KV
param msiKeyVaultName = '{{ .msiKeyVault.name }}'

// MGMT KV
param mgmtKeyVaultName = '{{ .mgmtKeyVault.name }}'

// MI for deployment scripts
param aroDevopsMsiId = '{{ .aroDevopsMsiId }}'

// Azure Monitor Workspace
param azureMonitorWorkspaceName = '{{ .monitoring.workspaceName }}'
