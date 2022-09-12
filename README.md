# secrets

A wrapper binary for storing simple key/value type secrets in Hashicorp Vault

## usage

Running SECRETS without any parameters shows the main top level help:

```
Usage:
  secrets [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage local configuration
  get         retrieve a secret
  delete      delete a secret  
  help        Help about any command
  list        list available secrets
  set         sets a secret

Flags:
  -h, --help      help for SECRETS
  -v, --version   version for SECRETS
```

Use the `-h` or `--help` flag to display help on any of the commands (and sub-commands)

## configuration

In order to use SECRETS, you need to ensure that you initially set the configuration - you can display the info needed by using the set help command:

```
Usage:
  secrets config set [flags]

Flags:
  -a, --address string    the address of Vault, e.g. http://127.0.0.1:9000
  -h, --help              help for set
  -w, --password string   the Vault password for the username
  -p, --project string    a project name (without spaces)
  -s, --store string      the Key Value store in Vault to use
  -u, --username string   the Vault username
```

Use the information sent to you by your a Vault administrator prior to trying to get/set any secrets.

## logical to actual storage

Each project has a list of secrets - access to these secrets (and the ability to read, write or delete) are controlled by the underlying Vault policy applied to the user. SECRETS makes the presumption that you have authority to do everything and then fails (gracefully!) if you don't.

- logical `store` is actually a Vault key/value store, i.e. a secret store off the root
- logical `project` is actually a path off the above secret store, e.g.  `/store/project`
- logical `secret` is a key/value pair with the key simply called `value` associated with the secret itself, e.g. `/store/project/app1-admin-password` could contain a single key/value pair called `value:myapp1password`
