# Vault Secret Tool
`kube-svc-ctl` is a command line tool for generating Helm 
compatible Yaml configuration files based on values stored in
Hashicorp Vault instance

## Installation

`kube-svc-ctl` is available as a downloadable binary from the [Nexus Tools repository](https://nexus.flotech.co/service/rest/repository/browse/tools/kube-svc-ctl/).

## Configuration
`kube-svc-ctl` respects `VAULT_ADDR` and `VAULT_TOKEN` 
environment variables.

Alternatively, `-vault-token` and `-vault-url` arguments can be used, but this approach is not 
recommended in pipelines due to possible credentials leak.

## Available commands
### generate-svc-config
### generate-helm-values
Generates Helm's YAML values file based on secret parameters in Hashicorp Vault.
Looks up `secret/data/service/$secret` secret in Vault where `$service` is the value of `-service` argument. 

#### Flags
`-service` Service name, used to lookup corresponding secrets in `secret/data/service/` directory in Vault.

`-add-docker-image-tag` Add information about service's docker tag to generated configuration. Default: `true`

`-tag` Service's docker image tag. Default: `latest`

`-target-yaml-tree` Target YAML tree in resulting values file. For example, specifying "some.target.tree" will result in:
```yaml
some:
  target:
    tree:
      EXAMPLE_KEY1: example_value1
      EXAMPLE_KEY2: example_value2
      ....
```
Default value: `secrets.datas`

`-target-yaml-format` How target values should be formatted. Supported formats: `list` (compatible with common way of specifying environment variables), `map`. Default: `map`

Example map formatting:
```yaml
some:
  target:
    tree:
      EXAMPLE_KEY1: example_value1
      EXAMPLE_KEY2: example_value2
      ....
```
Example list formatting:
```yaml
some:
  target:
    tree:
      - key: EXAMPLE_KEY1
        value: example_value1
      ....
```
```


### generate-secret-manifest
Generates a Kubernetes `Secret` of `kubernetes.io/dockerconfigjson` type based on values
in `secret/data/common/registry`  to authenticate with a container registry to pull a private image. 

Expected Vault directory secrets:

`registry_addr` docker registry address. Example: `registry.gitlab.com`

`registry_username` registry username

`registry_password` registry password
#### Flags
`-vault-url` Hashicorp Vault Url. Example: `http://localhost:8200`

`-vault-token` Hashicorp Vault Access Token
#### Usage examples
Generate and apply `Secret` with docker registry authentication credentials:

`kube-svc-ctl generate-secret-manifest | kubectl apply `

Example `Secret` manifest that will be generated:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: myregistrykey
data:
  .dockerconfigjson: UmVhbGx5IHJlYWxseSByZWVlZWVlZWVlZWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGx5eXl5eXl5eXl5eXl5eXl5eXl5eSBsbGxsbGxsbGxsbGxsbG9vb29vb29vb29vb29vb29vb29vb29vb29vb25ubm5ubm5ubm5ubm5ubm5ubm5ubm5ubmdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2cgYXV0aCBrZXlzCg==
type: kubernetes.io/dockerconfigjson
    
```