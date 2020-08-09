# Okay

Yet another Kubernetes manifests validation tool


## Usage

```
$> okay --help

A Kubernetes manifests validator tool

Usage:
  okay [flags]

Flags:
  -c, --crd stringArray        files or directories that contain CustomResourceDefinitions to be used for validation
  -f, --filename stringArray   (required) file or directories that contain the configuration to be validated
  -h, --help                   help for okay
  -R, --recursive              process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
  -s, --schema string          (required) Kubernetes OpenAPI V2 schema to validate against

```


![Image](img/results.png)
