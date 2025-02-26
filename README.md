# Salvo

A CLI for getting Kubernetes logs fast.

## Requirements

- Go 1.23+
- Kubernetes

## Installation

```shell
cd salvo/
go build .
go install
```

## Usage

The application comes with a few options for use. By default, running the application will use the default namespace of the current `kubeconfig` setup on your computer:

```sh
salvo logs # Default namespace, default output location of the current directory
```

A namespace can be specified instead of using `default`:

```sh
salvo logs --namespace backend # backend namespace, default output location of the current directory
salvo logs -n backend
```

More information is available via:

```sh
salvo -h
```
