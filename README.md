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

Before using the CLI, you need to authenticate:

```bash
nsai auth --token <your-token>
```

## Commands

### Initialize

Initialize a new project and cloud environment:

```bash
nsai init --cloud <provider> --project <project-name>
```

### Create Resources

Create various resources in Nstream AI:

```bash
nsai create --type <resource-type> --name <resource-name>
```

Resource types include:
- cluster
- provider
- tool

### Delete Resources

Delete existing resources:

```bash
nsai delete --type <resource-type> --name <resource-name>
```

### Get Information

Retrieve information about resources:

```bash
nsai get --type <resource-type> --name <resource-name>
```

### Update Resources

Update specific fields of resources:

```bash
nsai patch --type <resource-type> --name <resource-name>
```

### Set Active Context

Set the active project and cluster:

```bash
nsai use --project <project-name> --cluster <cluster-name>
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