# docker-credential-vault
Store and retrieve your docker registry credentials using Hashicorp Vault

Conforms to the docker credential helper API. 

Takes in a URL of the docker repository and returns a JSON object with populated credentials for the repository.

## Vault configuration

Vault configuration can be provided in three different ways. 

### Command line

*-vault*  hostname for vault server

*-port*   port used by vault server

*-token*  vault token with access to secrets for docker

### Configuration file

File format can be JSON, YAML, HCL, or TOML. Configuration should include values for 

*token*

*vault*

*port*

Default locations for the configuration file are 

 /etc/docker/vault-credential-helper/config.yaml
 
 $HOME/.docker/vault-credential-helper/config.yaml

The config file can also be specified by the -config command line option.

### Environment variables

*VAULT*  - hostname for vault

*TOKEN*  - token for vault to access credentials

*PORT*   - port for vault server
