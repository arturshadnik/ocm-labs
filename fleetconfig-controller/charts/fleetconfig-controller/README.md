# fleetconfig-controller helm chart

## TL;DR

```bash
helm repo add ocm https://open-cluster-management.io/helm-charts
helm repo update ocm
helm install fleetconfig-controller ocm/fleetconfig-controller -n fleetconfig-system --create-namespace
```

## Prerequisites

- Kubernetes >= v1.19
  
## Parameters

### FleetConfig Configuration

Configuration for the FleetConfig resource created on the Hub. By default, bootstraps the Hub cluster in hub-as-spoke mode.
### Spoke Feature Gates
Uncomment and configure `fleetConfig.spokeFeatureGates` to enable feature gates for the Klusterlet on each Spoke.
Do not disable the feature gates that are enabled by default.

Available Spoke Feature Gates:
- **AddonManagement** (ALPHA - default=true) - Enables addon management functionality
- **AllAlpha** (ALPHA - default=false) - Enables all alpha features
- **AllBeta** (BETA - default=false) - Enables all beta features
- **ClusterClaim** (ALPHA - default=true) - Enables cluster claim functionality
- **ExecutorValidatingCaches** (ALPHA - default=false) - Enables executor validating caches
- **RawFeedbackJsonString** (ALPHA - default=false) - Enables raw feedback JSON string support
- **V1beta1CSRAPICompatibility** (ALPHA - default=false) - Enables v1beta1 CSR API compatibility
### Registration Authentication Configuration
Registration authentication configuration for FleetConfig setup. Uses certificate signing requests (CSR) by default.
Available fields:
- **driver**: The authentication driver to use (default: "csr"). Set to "awsirsa" to use AWS IAM Roles for Service Accounts (IRSA) for EKS FleetConfigs.
- **hubClusterARN**: The ARN of the hub cluster (required for EKS FleetConfigs).
- **autoApprovedARNPatterns**: Optional list of spoke cluster ARN patterns that the hub will auto approve.
### Hub Cluster Manager Feature Gates
Feature gates for the Hub's Cluster Manager. Do not disable the feature gates that are enabled by default.

Available Hub Cluster Manager Feature Gates:
- **AddonManagement** (ALPHA - default=true) - Enables addon management functionality
- **AllAlpha** (ALPHA - default=false) - Enables all alpha features
- **AllBeta** (BETA - default=false) - Enables all beta features
- **CloudEventsDrivers** (ALPHA - default=false) - Enables cloud events drivers
- **DefaultClusterSet** (ALPHA - default=false) - Enables default cluster set functionality
- **ManagedClusterAutoApproval** (ALPHA - default=false) - Enables automatic managed cluster approval
- **ManifestWorkReplicaSet** (ALPHA - default=false) - Enables manifest work replica set functionality
- **NilExecutorValidating** (ALPHA - default=false) - Enables nil executor validation
- **ResourceCleanup** (BETA - default=true) - Enables automatic resource cleanup
- **V1beta1CSRAPICompatibility** (ALPHA - default=false) - Enables v1beta1 CSR API compatibility

### Singleton Control Plane Configuration
If provided, deploy a singleton control plane instead of Cluster Manager.
To enable singleton mode, `fleetConfig.hub.clusterManager` must be commented out and `fleetConfig.hub.singletonControlPlane` must be uncommented and configured with the following options:
- **name**: The name of the singleton control plane (default: "singleton-controlplane")
- **helm**: Helm configuration for the multicluster-controlplane Helm chart
  - **values**: Raw, YAML-formatted Helm values
  - **set**: List of comma-separated Helm values (e.g., key1=val1,key2=val2)
  - **setJson**: List of comma-separated Helm JSON values
  - **setLiteral**: List of comma-separated Helm literal STRING values
  - **setString**: List of comma-separated Helm STRING values

Refer to the [Multicluster Controlplane configuration](https://github.com/open-cluster-management-io/multicluster-controlplane/blob/main/charts/multicluster-controlplane/values.yaml) for more details.

| Name                                                                  | Description                                                                                                                                                                                                                                                  | Value                             |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------- |
| `fleetConfig.enabled`                                                 | Whether to create a FleetConfig resource.                                                                                                                                                                                                                    | `true`                            |
| `fleetConfig.spokeFeatureGates.ClusterClaim`                          | ClusterClaim feature gate (ALPHA - default=true). Enables cluster claim functionality.                                                                                                                                                                       | `true`                            |
| `fleetConfig.spokeFeatureGates.RawFeedbackJsonString`                 | RawFeedbackJsonString feature gate (ALPHA - default=false). Enables raw feedback JSON string support.                                                                                                                                                        | `true`                            |
| `fleetConfig.source.bundleVersion`                                    | Bundle version.                                                                                                                                                                                                                                              | `v1.0.0`                          |
| `fleetConfig.source.registry`                                         | Image registry.                                                                                                                                                                                                                                              | `quay.io/open-cluster-management` |
| `fleetConfig.registrationAuth.driver`                                 | The authentication driver to use (default: "csr"). Set to "awsirsa" to use AWS IAM Roles for Service Accounts (IRSA) for EKS FleetConfigs.                                                                                                                   | `csr`                             |
| `fleetConfig.hub.clusterManager.featureGates.DefaultClusterSet`       | DefaultClusterSet feature gate.                                                                                                                                                                                                                              | `true`                            |
| `fleetConfig.hub.clusterManager.featureGates.ManifestWorkReplicaSet`  | ManifestWorkReplicaSet feature gate.                                                                                                                                                                                                                         | `true`                            |
| `fleetConfig.hub.clusterManager.featureGates.ResourceCleanup`         | ResourceCleanup feature gate.                                                                                                                                                                                                                                | `true`                            |
| `fleetConfig.hub.clusterManager.purgeOperator`                        | If set, the cluster manager operator will be purged and the open-cluster-management namespace deleted when the FleetConfig CR is deleted.                                                                                                                    | `true`                            |
| `fleetConfig.hub.clusterManager.resources`                            | Resource specifications for all clustermanager-managed containers.                                                                                                                                                                                           | `{}`                              |
| `fleetConfig.hub.createNamespace`                                     | If true, create open-cluster-management namespace, otherwise use existing one.                                                                                                                                                                               | `true`                            |
| `fleetConfig.hub.force`                                               | If set, the hub will be reinitialized.                                                                                                                                                                                                                       | `false`                           |
| `fleetConfig.hub.kubeconfig.context`                                  | The context to use in the kubeconfig file. Leave empty to use the current context.                                                                                                                                                                           | `""`                              |
| `fleetConfig.hub.kubeconfig.inCluster`                                | If set, the kubeconfig will be read from the cluster. Only applicable for same-cluster operations.                                                                                                                                                           | `true`                            |
| `fleetConfig.spokes[0].name`                                          | Name of the spoke cluster.                                                                                                                                                                                                                                   | `hub-as-spoke`                    |
| `fleetConfig.spokes[0].createNamespace`                               | If true, create open-cluster-management namespace and agent namespace (open-cluster-management-agent for Default mode, <klusterlet-name> for Hosted mode), otherwise use existing one. Do not edit this name if you are using the default hub-as-spoke mode. | `true`                            |
| `fleetConfig.spokes[0].syncLabels`                                    | If true, sync the labels from klusterlet to all agent resources.                                                                                                                                                                                             | `false`                           |
| `fleetConfig.spokes[0].kubeconfig.context`                            | The context to use in the kubeconfig file. Leave empty to use the current context.                                                                                                                                                                           | `""`                              |
| `fleetConfig.spokes[0].kubeconfig.inCluster`                          | If set, the kubeconfig will be read from the cluster. Only applicable for same-cluster operations.                                                                                                                                                           | `true`                            |
| `fleetConfig.spokes[0].ca`                                            | Hub cluster CA certificate, optional.                                                                                                                                                                                                                        | `""`                              |
| `fleetConfig.spokes[0].proxyCa`                                       | Proxy CA certificate, optional.                                                                                                                                                                                                                              | `""`                              |
| `fleetConfig.spokes[0].proxyUrl`                                      | URL of a forward proxy server used by agents to connect to the Hub cluster, optional.                                                                                                                                                                        | `""`                              |
| `fleetConfig.spokes[0].klusterlet.mode`                               | Deployment mode for klusterlet. Options: Default (agents run on spoke cluster) | Hosted (agents run on hub cluster).                                                                                                                                         | `Default`                         |
| `fleetConfig.spokes[0].klusterlet.purgeOperator`                      | If set, the klusterlet operator will be purged and all open-cluster-management namespaces deleted when the klusterlet is unjoined from its Hub cluster.                                                                                                      | `true`                            |
| `fleetConfig.spokes[0].klusterlet.forceInternalEndpointLookup`        | If true, the klusterlet agent will start the cluster registration process by looking for the                                                                                                                                                                 | `true`                            |
| `fleetConfig.spokes[0].klusterlet.managedClusterKubeconfig`           | External managed cluster kubeconfig, required if using hosted mode.                                                                                                                                                                                          | `{}`                              |
| `fleetConfig.spokes[0].klusterlet.forceInternalEndpointLookupManaged` | If true, the klusterlet accesses the managed cluster using the internal endpoint from the public cluster-info in the managed cluster instead of using managedClusterKubeconfig.                                                                              | `false`                           |
| `fleetConfig.spokes[0].klusterlet.resources`                          | Resource specifications for all klusterlet-managed containers.                                                                                                                                                                                               | `{}`                              |
| `fleetConfig.spokes[0].klusterlet.singleton`                          | If true, deploy klusterlet in singleton mode, with registration and work agents running in a single pod. This is an alpha stage flag.                                                                                                                        | `false`                           |

### Topology Resources

| Name                        | Description                                                                                                                                                                                                                                                                                                                                        | Value  |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| `topologyResources.enabled` | Whether to create Namespaces and ManagedClusterSetBindings for the default ManagedClusterSets created when a FleetConfig is created with the DefaultClusterSet feature gate enabled. Additionally, a Namespace, ManagedClusterSet, and Placement are created for targeting all managed clusters that are not the hub running in hub-as-spoke mode. | `true` |

### fleetconfig-controller parameters

| Name                                                | Description                                                                                                                            | Value                                                    |
| --------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| `kubernetesProvider`                                | Kubernetes provider of the cluster that fleetconfig-controller will be installed on. Valid values are "Generic", "EKS", "GKE-Ingress". | `Generic`                                                |
| `replicas`                                          | fleetconfig-controller replica count                                                                                                   | `1`                                                      |
| `imageRegistry`                                     | Image registry                                                                                                                         | `""`                                                     |
| `image.repository`                                  | Image repository                                                                                                                       | `quay.io/open-cluster-management/fleetconfig-controller` |
| `image.tag`                                         | Image tag                                                                                                                              | `v0.0.5`                                                 |
| `image.pullPolicy`                                  | Image pull policy                                                                                                                      | `IfNotPresent`                                           |
| `imagePullSecrets`                                  | Image pull secrets                                                                                                                     | `[]`                                                     |
| `serviceAccount.annotations`                        | Annotations to add to the service account                                                                                              | `{}`                                                     |
| `containerSecurityContext.allowPrivilegeEscalation` | allowPrivilegeEscalation                                                                                                               | `false`                                                  |
| `containerSecurityContext.capabilities.drop`        | capabilities to drop                                                                                                                   | `["ALL"]`                                                |
| `containerSecurityContext.runAsNonRoot`             | runAsNonRoot                                                                                                                           | `true`                                                   |
| `resources.limits.cpu`                              | fleetconfig controller's cpu limit                                                                                                     | `500m`                                                   |
| `resources.limits.memory`                           | fleetconfig controller's memory limit                                                                                                  | `256Mi`                                                  |
| `resources.requests.cpu`                            | fleetconfig controller's cpu request                                                                                                   | `100m`                                                   |
| `resources.requests.memory`                         | fleetconfig controller's memory request                                                                                                | `256Mi`                                                  |
| `healthCheck.port`                                  | port the liveness & readiness probes are bound to                                                                                      | `9440`                                                   |
| `kubernetesClusterDomain`                           | kubernetes cluster domain                                                                                                              | `cluster.local`                                          |

### cert-manager

| Name                   | Description                      | Value  |
| ---------------------- | -------------------------------- | ------ |
| `cert-manager.enabled` | Whether to install cert-manager. | `true` |

### certificates

| Name                                         | Description                                 | Value                    |
| -------------------------------------------- | ------------------------------------------- | ------------------------ |
| `certificates.clusterIssuer.spec.selfSigned` | Use a self-signed ClusterIssuer by default. | `{}`                     |
| `certificates.clusterIssuer.enabled`         | Enable the creation of a ClusterIssuer.     | `true`                   |
| `certificates.issuerRef.kind`                | Kind of the certificate issuer to use.      | `ClusterIssuer`          |
| `certificates.issuerRef.name`                | Name of the certificate issuer to use.      | `fleetconfig-controller` |

### webhook parameters

| Name                                                 | Description                              | Value                    |
| ---------------------------------------------------- | ---------------------------------------- | ------------------------ |
| `admissionWebhooks.enabled`                          | enable admission webhooks                | `true`                   |
| `admissionWebhooks.failurePolicy`                    | admission webhook failure policy         | `Fail`                   |
| `admissionWebhooks.certificate.mountPath`            | admission webhook certificate mount path | `/etc/k8s-webhook-certs` |
| `admissionWebhooks.certManager.revisionHistoryLimit` | cert-manager revision history limit      | `3`                      |
| `webhookService.type`                                | webhook service type                     | `ClusterIP`              |
| `webhookService.port`                                | webhook service port                     | `9443`                   |

### dev parameters

| Name               | Description       | Value                    |
| ------------------ | ----------------- | ------------------------ |
| `devspaceEnabled`  | devspace enabled  | `false`                  |
| `fullnameOverride` | Fullname override | `fleetconfig-controller` |
