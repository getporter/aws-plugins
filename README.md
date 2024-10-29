# AWS Plugins for Porter

This is a set of AWS plugins for [Porter](https://github.com/getporter/porter).

## Install

The plugin is distributed as a single binary, `aws`. The following snippet will clone this repository, build the binary
and install it to **~/.porter/plugins/**.

```
go get get.porter.sh/plugin/aws/cmd/aws
cd $(go env GOPATH)/src/get.porter.sh/plugin/aws
mage build install
```

After installing the plugin, you must modify your porter configuration file and select which plugin you want to use.

## Secrets

Secrets plugins allow Porter to inject secrets into credential or parameter sets. It also stores sensitive data referenced/generated during Porter execution.

For example, if your team has a shared key vault with a database password, you
can use the secrets manager plugin to inject it as a credential or parameter when you install a bundle.

### Secrets Manager

The `aws.secretsmanager` plugin resolves credentials or parameters against secrets in AWS Secrets Manager. It's also used to store any sensitive data referenced during Porter execution.

1. Open, or create, `~/.porter/config.toml`
1. Add the following lines to activate the AWS Secrets Manager secrets plugin:

    ```toml
    default-secrets = "mysecrets"
    
    [[secrets]]
    name = "mysecrets"
    plugin = "aws.secretsmanager"
    
    [secrets.config]
    region = "us-east-2"
    ```
1. [Create a Secrets Manager][secretsmanager] and set the region in the config.

### Secret ID
Both the secret name and secret ARN can be used to resolve a secret that may not exist in the plugin configured Secrets Manager.

An example CredentialSet that would resolve to both the configured AWS Secrets Manager as well as a separate AWS Secrets Manager based on the ARN would look like this:

```yaml
name: example-credset
schemaVersion: 1.0.1
credentials:
  - name: example-configured-secret
    source:
      secret: my-secret
  - name: example-secret-id
    source:
      secret: arn:aws:secretsmanager:us-east-2:111122223333:secret:my-secret-abcdef
```

This provides `porter` with the ability to fetch secrets out of multiple AWS Secrets Managers without having the change the default configuration. 

### Authentication

Authentication to AWS can be done in different ways, e.g., using environment variables or AWS credentials file, see the AWS documentation,
[Specifying Credentials][specify-creds], for details.

[secretsmanager]: https://aws.amazon.com/secrets-manager/
[specify-creds]: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/#specifying-credentials
