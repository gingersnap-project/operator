# engytita-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

We utilise [Skaffold](https://skaffold.dev/) to drive CI/CD, so you will need to download the latest binary in order to
follow the steps below:

### Kind Cluster

Create a local kind cluster backed by a local docker repository, with [OLM](https://olm.operatorframework.io/) and
[cert-manager](https://cert-manager.io) installed:

```sh
./hack/kind.sh`
```

### Development

Build the Operator image and deploy to a cluster:

```sh
skaffold dev
```

Changes to the local `**/*.go` files will result in the image being rebuilt and the Operator deployment updated. 

### Debugging
Build the Operator image with [dlv](https://github.com/go-delve/delve) so that a remote debugger can be attached
to the Operator deployment from your IDE.

```sh
skaffold debug
```

### Deploying
Build the Operator image and deploy to a cluster:

```sh
skaffold run
```

### Remote Repositories
The `skaffold dev|debug|run` commands can all be used on a remote k8s instance, as long as the built images are accessible
on the cluster. To build and push the operator images to a remote repository, add the `--default-repo` option, for example:

```sh
skaffold run --default-repo <remote_repo>
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

### Testing
The project consists of three distinct types of test:

1. unit
2. integration
3. e2e

These tests are based upon the [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/) libraries.

#### Unit
Unit tests should be created for packages using the go `_test.go` convention.

#### Integration
Controller and webhook integration tests are implemented using [envtest](https://book.kubebuilder.io/cronjob-tutorial/writing-tests.html).
These tests should be executable using only the configured controller/webhooks and the local k8s api-server. They
shouldn't rely on other k8 controllers, such as the pod Controller.

#### E2E
Located in `test/e2e` dir. These tests also utilise`envtest` but rely on an existing k8s cluster to provide controllers
for core components, such as Pods. Any test that depends on a controller defined outside this project should be
implemented as an E2E test.

All E2E tests should be annotated with the following build tags to ensure that they are only executed with `make test-e2e`:

```go
//go:build e2e
// +build e2e
```

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
