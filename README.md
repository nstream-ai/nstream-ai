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
1. Selecting cluster type (Basic/Standard/Enterprise)
2. Choosing cloud provider (AWS/GCP/Azure)
3. Selecting region for the chosen cloud provider
4. Setting up bucket access:
   - Shows existing buckets compatible with the selected cloud provider
   - Option to use an existing bucket or create a new one
   - If creating new, prompts for bucket name
5. Setting up bucket access role/principal
6. Verifying bucket access
7. Creating the cluster

Example:
```bash
# Interactive mode
nsai create cluster my-cluster

# With flags
nsai create cluster my-cluster --type basic --cloud aws --region us-east-1 --bucket my-bucket --role my-role
```

#### Create Bucket

```bash
nsai create bucket [bucket-name] [flags]
```

Flags:
- `--name, -n`: Bucket name (optional, will prompt if not provided)
- `--provider, -p`: Cloud provider (aws/gcp/azure)
- `--region, -r`: Region for the bucket

The command will:
1. Check for existing buckets compatible with the cluster's cloud provider
2. If compatible buckets exist:
   - Display them in a table with details
   - Ask if you want to use an existing bucket
   - If yes, let you select one from the list
   - If no, proceed with creating a new bucket
3. If creating a new bucket:
   - Prompt for bucket name (if not provided)
   - Select region based on cloud provider
   - Create the bucket

Example:
```bash
# Interactive mode
nsai create bucket

# With name
nsai create bucket my-bucket

# With all flags
nsai create bucket my-bucket --provider aws --region us-east-1
```

### Use Resources

#### Use Cluster

```bash
nsai use cluster [cluster-name] [flags]
```

Flags:
- `--name, -n`: Cluster name (optional, will prompt for selection if not provided)

The command will:
1. List available clusters
2. Let you select a cluster to use
3. Set the selected cluster as the current context

Example:
```bash
# Interactive mode
nsai use cluster

# Direct selection
nsai use cluster my-cluster
```

#### Use Bucket

```bash
nsai use bucket [bucket-name] [flags]
```

Flags:
- `--name, -n`: Bucket name (optional, will prompt for selection if not provided)
- `--cluster, -c`: Cluster name (optional, will use current cluster or prompt for selection)

The command will:
1. Check if a cluster context is set or provided
2. List available buckets compatible with the cluster's cloud provider
3. Let you select a bucket to use
4. Verify cloud provider compatibility between cluster and bucket
5. Set the selected bucket as the current context

Example:
```bash
# Interactive mode
nsai use bucket

# With bucket name
nsai use bucket my-bucket

# With cluster context
nsai use bucket my-bucket --cluster my-cluster
```

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