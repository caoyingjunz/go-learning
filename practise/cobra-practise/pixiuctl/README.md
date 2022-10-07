# Pixiuctl Usage Overview

## Global Help
```go
# go run pixiuctl.go help
pixiuctl controls the Pixiu cluster manager.

 Find more information at: https://github.com/caoyingjunz/go-learning

Basic Commands (Beginner):
  create      Create a pixiu resource from a file or stdin

Advanced Commands:
  apply       Apply a configuration to a pixiu resource by file name or stdin

Other Commands:
  completion  generate the autocompletion script for the specified shell

Usage:
  pixiuctl [flags] [options]

Use "pixiuctl <command> --help" for more information about a given command.
```

## Subcommand Help
```go
# go run pixiuctl.go create --help
Create a resource from a file or from stdin.

 JSON and YAML formats are accepted.

Examples:
  pixiuctl create -f ./create.json

Available Commands:
  service     Create a ClusterIP service

Options:
      --edit=false: Edit the API resource before creating
      --raw='': Raw URI to POST to the server.  Uses the transport specified by the kubeconfig file.

Usage:
  pixiuctl create -f FILENAME [options]

Use "pixiuctl <command> --help" for more information about a given command.
Use "pixiuctl create options" for a list of global command-line options (applies to all commands).
```