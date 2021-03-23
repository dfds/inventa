<h1 align="center">Welcome to inventa-operator 👋</h1>
<p>
  <a href="api-docs.md" target="_blank">
    <img alt="Documentation" src="https://img.shields.io/badge/documentation-yes-brightgreen.svg" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> Operator that monitors Kubernetes objects and exposes them through an API

## Install

The operator is available as a [Helm](https://helm.sh/) chart. Chart files can be found here [https://github.com/dfds/helm-charts](https://github.com/dfds/helm-charts/tree/main/charts/inventa-operator).

```sh
helm repo add dfds https://dfds.github.io/helm-charts
helm repo update
helm install inven dfds/inventa-operator
```

Running the commands above will install inventa-operator with the default settings. If you wish you customise beyond stock, you can find all the support configuration values in the [values.yaml](https://github.com/dfds/helm-charts/blob/main/charts/inventa-operator/values.yaml). The core options are the following:

```yaml
inventa:
  # Enables monitoring of ServiceProxyCrd
  enableServiceProxyCrd: false
  # Enable Ingress annotation controller
  enableIngressMonitoring: true
  # Enable Service annotation controller
  enableServiceMonitoring: true  
  # Enable HTTP API
  enableHttpApi: true
  # Enables auth on all HTTP endpoints(not including metrics). Currently only supports Azure AD. authClientId and authTenantId must be configured if this is enabled.
  enableHttpApiAuth: false
  # Specify IPv4 address to bind the API server to.
  bindAddress: "127.0.0.1"  
  # Azure AD client id. Not necessary if enableHttpApiAuth is set to false
  authClientId: ""
  # Azure AD tenant id. Not necessary if enableHttpApiAuth is set to false 
  authTenantId: ""

rbac:
  # Will add RBAC manifests(clusterrole, clusterrolebinding, serviceaccount) as a part of the Chart installation
  create: true
  # Will ONLY add RBAC manifests. It won't deploy the operator
  installOnlyRbac: false  
```

### RBAC

The `rbac.installOnlyRbac` option exists in case you deploy/update Helm installations in a pipeline with non-elevated rights(e.g. not allowed to create clusterroles, clusterrolebindings, etc..). With that option, one could run `helm template inven dfds/inventa-operator --set rbac.installOnlyRbac=true` and generate the RBAC manifests separately from the Chart installation.

## Usage

With the operator running with the default settings, the following will be running:

- Monitoring of Ingress and Service objects cluster-wide
- HTTP API

Documentation for the API can be found in [openapi.yaml](operator/openapi.yaml).

## Author

👤 **DFDS A/S: Richard Fisher [@rifisdfds](https://github.com/rifisdfds), Emil H. Clausen [@SEQUOIIA](https://github.com/SEQUOIIA)**

* Website: https://dfds.com/en/tech
* Github: [@dfds](https://github.com/dfds)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/dfds/inventa/issues?q=label%3Aoperator).
