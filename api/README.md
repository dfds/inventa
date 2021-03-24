<h1 align="center">inventa-api</h1>
<p>
  <a href="api-docs.md" target="_blank">
    <img alt="Documentation" src="https://img.shields.io/badge/documentation-yes-brightgreen.svg" />
  </a>
  <a href="#" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>
</p>

> A service that can connect to one or more inventa-operator instances, aggregate results from the operator(s) and make them available through an API

## Install

The api is available as a [Helm](https://helm.sh/) chart. Chart files can be found here [https://github.com/dfds/helm-charts](https://github.com/dfds/helm-charts/tree/main/charts/inventa-api).

```sh
helm repo add dfds https://dfds.github.io/helm-charts
helm repo update
helm install inven-api dfds/inventa-api
```

Running the commands above will install inventa-api with the default settings. If you wish to customise beyond stock, you can find all the support configuration values in the [values.yaml](https://github.com/dfds/helm-charts/blob/main/charts/inventa-api/values.yaml).

## Usage

With the operator running with the default settings, the following will be running:

- HTTP API

Documentation for the API can be found in [openapi.yaml](api/openapi.yaml). An always up-to-date *openapi.yaml* can also be found by running the application and going to `/swagger/v1/swagger.yaml`.

## Author

👤 **DFDS A/S: Richard Fisher [@rifisdfds](https://github.com/rifisdfds), Emil H. Clausen [@SEQUOIIA](https://github.com/SEQUOIIA)**

* Website: https://dfds.com/en/tech
* Github: [@dfds](https://github.com/dfds)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/dfds/inventa/issues?q=label%3Aapi).
