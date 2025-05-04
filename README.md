# Nstream AI CLI

A command-line interface for managing Nstream AI resources and services.

## Installation

### Using Go

```bash
go install github.com/nstream-ai/nstream-ai-cli@latest
```

### Using Homebrew (macOS/Linux)

```bash
brew tap nstream-ai/nstream-ai-cli
brew install nstream-ai-cli
```

## Authentication

Before using the CLI, you need to authenticate. The CLI provides two authentication methods:

### Sign In

```bash
nsai auth signin
```

This will guide you through an interactive sign-in process where you'll:
1. Enter your email
2. Receive a password via email
3. Enter the password to complete authentication

### Sign Up

```bash
nsai auth signup
```

This will guide you through creating a new account where you'll:
1. Enter your email
2. Enter your organization name
3. Enter your name
4. Enter your role
5. Receive a password via email
6. Enter the password to complete registration

## Commands

### Initialize

Initialize a new project and cloud environment:

```bash
nsai init [flags]
```

Flags:
- `--user, -u`: Username for authentication
- `--password, -p`: Password for authentication
- `--cluster, -c`: Cluster name to use
- `--create-cluster`: Create a new cluster (requires additional flags)
- `--region`: Region for cluster creation
- `--cloud`: Cloud provider for cluster creation
- `--bucket`: Bucket name for cluster creation
- `--role`: Role for cluster creation

### Create Resources

#### Create Cluster

```bash
nsai create cluster <cluster-name> [flags]
```

Flags:
- `--type, -t`: Cluster type (basic/standard/enterprise) [default: basic]
- `--cloud, -c`: Cloud provider (aws/gcp/azure) [default: gcp]
- `--region, -r`: Region for the cluster
- `--bucket, -b`: Bucket name for storage
- `--role, -p`: Role/principal to assume for bucket access

The command will guide you through:
1. Selecting cluster type
2. Choosing cloud provider
3. Selecting region
4. Setting up bucket access
5. Verifying bucket access
6. Creating the cluster

### Delete Resources

```bash
nsai delete [resource-type] [resource-name]
```

### Get Information

```bash
nsai get [resource-type] [resource-name]
```

### Update Resources

```bash
nsai patch [resource-type] [resource-name]
```

### Set Active Context

```bash
nsai use [flags]
```

Flags:
- `--project, -p`: Project to set as active
- `--cluster, -c`: Cluster to set as active

## Global Flags

- `-v, --verbose`: Enable verbose output
- `-h, --help`: Show help for command

## Shell Completion

The CLI supports shell completion for Bash, Zsh, and Fish. To enable:

### Bash

```bash
source <(nsai completion bash)
```

### Zsh

```bash
source <(nsai completion zsh)
```

### Fish

```bash
source (nsai completion fish | psub)
```

## License

MIT License - see [LICENSE](LICENSE) for details