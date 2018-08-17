# vault-gcp-authenticator

The `vault-gcp-authenticator` is a small application/container that performs the [HashiCorp Vault][vault] [GCP authentication process][vault-gcp-auth]
and places the Vault token in a well-known, configurable location or prints to STDOUT.

[vault]: https://www.vaultproject.io
[vault-gcp-auth]: https://www.vaultproject.io/docs/auth/kubernetes.html#authentication

This project was forked from https://github.com/sethvargo/vault-kubernetes-authenticator.
The GCP and K8S Vault auth backends are very similar with the primary difference being the source
of the JWT token.

Key changes in this repo:

* Reads instance identity JWT from the Google Cloud Platform's HTTP metadata API. https://cloud.google.com/compute/docs/instances/verifying-instance-identity
* Uses Vault go client for Vault interactions instead of the Go http client. This allows for all of the standard Vault environment variables to be used.
* Supports writing Vault token to STDOUT for easier integration in scripts.
* Arguments can be supplied as flags in addition to environment vars. Vault client args must be specified in the environment.

## Configuration

- `-d, --destination, $TOKEN_DEST_PATH` - The path on disk to store the token, or "-" for stdout. (default: /.vault-token)

- `-r, --role, $VAULT_ROLE` - **Required** The name of the Vault GCP role to use for authentication

- `-m, --metadata-addr, $METADATA_ADDR` - Hostname or IP of the GCP metadata API. (default: metadata.google.internal). This can be useful if you do not use GCP's DNS servers in which case you can specify the IP: `169.254.169.254`

- `-p, --path, $VAULT_GCP_MOUNT_PATH` - The name of the mount where the GCP auth method is enabled. (default: gcp)

```text
vault auth enable -path=google gcp -> VAULT_GCP_MOUNT_PATH=google
```

- The Vault Go Library is used to manage Vault communications which means all of the standard Vault environment variables are available
  such as `VAULT_ADDR`, `VAULT_CAPATH`, `VAULT_CACERT`, etc. https://www.vaultproject.io/docs/commands/index.html#environment-variables

## Example Usage

```shell
#!/bin/bash

set -eou pipefail

export VAULT_ADDR="https://vault:8200"
export VAULT_TOKEN=$(/bin/vault-gcp-authenticator -r myapp-role)

vault read secret/foo >/etc/myapp/secrets
vault write pki/issue/cert common_name="myapp" # parse key and cert, store to disk
```
