# Okay

Yet another Kubernetes manifests validation tool


## Usage

```
$> okay --help

A Kubernetes manifests validator tool

Usage:
  okay [flags]

Flags:
  -f, --filename stringArray   that contains the configuration to be validated (default [.])
  -h, --help                   help for okay
  -R, --recursive              Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
  -s, --schema string          Kubernetes OpenAPI schema
```
