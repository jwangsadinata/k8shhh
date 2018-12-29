# k8shhh

<img src="https://i.imgur.com/OQN6tuT.png" width="400">

Quickly encode your configuration into K8s secrets.

### Contents

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [Contributing](#contributing)

## Features

- Written in simple [Go][go-project]
- No installation necessary - binary is provided
- Intuitive and [easy to use][usage]
- Supports encoding to both `JSON` and `YAML`
- Works on Linux, Mac and Windows

## Installation

### Option 1: Binary

Download the latest binary from the [Releases][releases] page. This is the
easiest way to get started with `k8shhh`.

Don't forget to add the location of the binary to your `$PATH`.

### Option 2: From source

Install the package via the following:

    go get -u github.com/jwangsadinata/k8shhh

## Usage

### Commands:

```
encode  | encode key value as a kubernetes secret
decode  | decode the kubernetes secret into a readable configuration
version | print the current version of k8shhh
```

### Some Examples:

#### Encode from the standard input

```bash
$ echo "ETCD_NAME=kube-etcd" | k8shhh encode
apiVersion: v1
data:
  ETCD_NAME: a3ViZS1ldGNk
kind: Secret
metadata:
  name: mysecret
type: Opaque

$ cat <<EOF | k8shhh encode
ETCD_NAME=kube-etcd
ETCD_SNAPSHOT_COUNT=5000
ETCD_HEARTBEAT_INTERVAL=100
ETCD_ELECTION_TIMEOUT=500
ETCD_LISTEN_PEER_URLS=http://10.0.0.1:2380
ETCD_LISTEN_CLIENT_URLS=http://10.0.0.1:2379
EOF
apiVersion: v1
data:
  ETCD_ELECTION_TIMEOUT: NTAw
  ETCD_HEARTBEAT_INTERVAL: MTAw
  ETCD_LISTEN_CLIENT_URLS: aHR0cDovLzEwLjAuMC4xOjIzNzk=
  ETCD_LISTEN_PEER_URLS: aHR0cDovLzEwLjAuMC4xOjIzODA=
  ETCD_NAME: a3ViZS1ldGNk
  ETCD_SNAPSHOT_COUNT: NTAwMA==
kind: Secret
metadata:
  name: mysecret
type: Opaque
```

The result can then be outputted into a `.yaml` or `.json` file, like the
following:

```bash
$ cat <<EOF | k8shhh encode -o example
ETCD_NAME=kube-etcd
ETCD_SNAPSHOT_COUNT=5000
ETCD_HEARTBEAT_INTERVAL=100
ETCD_ELECTION_TIMEOUT=500
ETCD_LISTEN_PEER_URLS=http://10.0.0.1:2380
ETCD_LISTEN_CLIENT_URLS=http://10.0.0.1:2379
EOF
example.yaml

$ cat example.yaml
apiVersion: v1
data:
  ETCD_ELECTION_TIMEOUT: NTAw
  ETCD_HEARTBEAT_INTERVAL: MTAw
  ETCD_LISTEN_CLIENT_URLS: aHR0cDovLzEwLjAuMC4xOjIzNzk=
  ETCD_LISTEN_PEER_URLS: aHR0cDovLzEwLjAuMC4xOjIzODA=
  ETCD_NAME: a3ViZS1ldGNk
  ETCD_SNAPSHOT_COUNT: NTAwMA==
kind: Secret
metadata:
  name: example
type: Opaque
```

```bash
$ cat <<EOF | k8shhh encode -f json -o example
ETCD_NAME=kube-etcd
ETCD_SNAPSHOT_COUNT=5000
ETCD_HEARTBEAT_INTERVAL=100
ETCD_ELECTION_TIMEOUT=500
ETCD_LISTEN_PEER_URLS=http://10.0.0.1:2380
ETCD_LISTEN_CLIENT_URLS=http://10.0.0.1:2379
EOF
example.json

$ cat example.json
{
    "apiVersion": "v1",
    "data": {
        "ETCD_ELECTION_TIMEOUT": "NTAw",
        "ETCD_HEARTBEAT_INTERVAL": "MTAw",
        "ETCD_LISTEN_CLIENT_URLS": "aHR0cDovLzEwLjAuMC4xOjIzNzk=",
        "ETCD_LISTEN_PEER_URLS": "aHR0cDovLzEwLjAuMC4xOjIzODA=",
        "ETCD_NAME": "a3ViZS1ldGNk",
        "ETCD_SNAPSHOT_COUNT": "NTAwMA=="
    },
    "kind": "Secret",
    "metadata": {
        "name": "example"
    },
    "type": "Opaque"
}
```

#### Encode from existing file

```bash
$ ls
example-file

$ cat example-file
DB_HOST=localhost
DB_PORT=5432

$ k8shhh encode -i example-file
apiVersion: v1
data:
  DB_HOST: bG9jYWxob3N0
  DB_PORT: NTQzMg==
kind: Secret
metadata:
  name: mysecret
type: Opaque
```

Similarly, the result can also be outputted into a proper `.yaml` or
`.json` file, as follows:

```bash
$ k8shhh encode -i example-file -o example
example.yaml

$ cat example.yaml
apiVersion: v1
data:
  DB_HOST: bG9jYWxob3N0
  DB_PORT: NTQzMg==
kind: Secret
metadata:
  name: example
type: Opaque
```

```bash
$ k8shhh encode -f json -n local-postgres -i example-file -o example
example.json

$ cat example.json
{
    "apiVersion": "v1",
    "data": {
        "DB_HOST": "bG9jYWxob3N0",
        "DB_PORT": "NTQzMg=="
    },
    "kind": "Secret",
    "metadata": {
        "name": "local-postgres"
    },
    "type": "Opaque"
}
```

#### Using kubectl with `k8shhh encode`

`k8shhh encode` works well with [kubectl][kubectl], which is the command line
client for kubernetes. Some of the examples include the following:

```bash
$ kubectl create -f $(echo "TOKEN=8fd41973acac04e5fc76fde5439c8b94f1eb1233" | k8shhh encode -o token)
secret "token" created

$ kubectl describe secrets/token
Name:         token
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
TOKEN:  40 bytes
```

The most common use case will be to encode directly from a relatively long
`.env` file, and pass it as kubernetes secrets, as follows:

```bash
$ ls -a
.env

$ cat .env
CI_JOB_ID="50"
CI_COMMIT_SHA="1ecfd275763eff1d6b4844ea3168962458c9f27a"
CI_COMMIT_REF_NAME="master"
CI_REPOSITORY_URL="https://gitlab-ci-token:abcde-1234ABCD5678ef@example.com/gitlab-org/gitlab-ce.git"
CI_COMMIT_TAG="1.0.0"
CI_JOB_NAME="spec:other"
CI_JOB_STAGE="test"
CI_JOB_MANUAL="true"
CI_JOB_TRIGGERED="true"
CI_JOB_TOKEN="abcde-1234ABCD5678ef"
CI_PIPELINE_ID="1000"
CI_PIPELINE_IID="10"
CI_PROJECT_ID="34"
CI_PROJECT_DIR="/builds/gitlab-org/gitlab-ce"
CI_PROJECT_NAME="gitlab-ce"
CI_PROJECT_NAMESPACE="gitlab-org"
CI_PROJECT_PATH="gitlab-org/gitlab-ce"
CI_PROJECT_URL="https://example.com/gitlab-org/gitlab-ce"
CI_REGISTRY="registry.example.com"
CI_REGISTRY_IMAGE="registry.example.com/gitlab-org/gitlab-ce"
CI_RUNNER_ID="10"
CI_RUNNER_DESCRIPTION="my runner"
CI_RUNNER_TAGS="docker, linux"
CI_SERVER="yes"
CI_SERVER_NAME="GitLab"
CI_SERVER_REVISION="70606bf"
CI_SERVER_VERSION="8.9.0"
CI_SERVER_VERSION_MAJOR="8"
CI_SERVER_VERSION_MINOR="9"
CI_SERVER_VERSION_PATCH="0"
GITLAB_USER_ID="42"
GITLAB_USER_EMAIL="user@example.com"
CI_REGISTRY_USER="gitlab-ci-token"
CI_REGISTRY_PASSWORD="longalfanumstring"

$ kubectl create -f $(k8shhh encode -i .env -o gitlab)
secret "gitlab" created

$ kubectl describe secrets/gitlab
Name:         gitlab
Namespace:    default
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
GITLAB_USER_EMAIL:        16 bytes
CI_SERVER_VERSION_MINOR:  1 bytes
CI_PROJECT_ID:            2 bytes
CI_REGISTRY:              20 bytes
CI_JOB_STAGE:             4 bytes
CI_PROJECT_PATH:          20 bytes
CI_PROJECT_URL:           40 bytes
CI_REPOSITORY_URL:        81 bytes
CI_COMMIT_REF_NAME:       6 bytes
CI_COMMIT_SHA:            40 bytes
CI_JOB_TRIGGERED:         4 bytes
CI_PIPELINE_IID:          2 bytes
CI_PROJECT_DIR:           28 bytes
CI_PROJECT_NAME:          9 bytes
CI_RUNNER_DESCRIPTION:    9 bytes
CI_SERVER:                3 bytes
CI_JOB_ID:                2 bytes
CI_JOB_TOKEN:             20 bytes
CI_SERVER_VERSION_MAJOR:  1 bytes
GITLAB_USER_ID:           2 bytes
CI_PIPELINE_ID:           4 bytes
CI_PROJECT_NAMESPACE:     10 bytes
CI_JOB_NAME:              10 bytes
CI_REGISTRY_PASSWORD:     17 bytes
CI_SERVER_REVISION:       7 bytes
CI_SERVER_VERSION:        5 bytes
CI_COMMIT_TAG:            5 bytes
CI_JOB_MANUAL:            4 bytes
CI_SERVER_VERSION_PATCH:  1 bytes
CI_RUNNER_TAGS:           13 bytes
CI_SERVER_NAME:           6 bytes
CI_RUNNER_ID:             2 bytes
CI_REGISTRY_IMAGE:        41 bytes
CI_REGISTRY_USER:         15 bytes
```

#### Decode from standard input

```bash
$ echo "{\"apiVersion\":\"v1\",\"data\":{\"DB_HOST\":\"bG9jYWxob3N0\",\"DB_PORT\":\"NTQzMg==\"},\"kind\":\"Secret\",\"metadata\":{\"name\":\"postgres\"},\"type\":\"Opaque\"}" | k8shhh decode
DB_HOST=localhost
DB_PORT=5432
```

The result can then be easily outputted into a file by using the following
command:

```bash
$ echo "{\"apiVersion\":\"v1\",\"data\":{\"DB_HOST\":\"bG9jYWxob3N0\",\"DB_PORT\":\"NTQzMg==\"},\"kind\":\"Secret\",\"metadata\":{\"name\":\"postgres\"},\"type\":\"Opaque\"}" | k8shhh decode -o postgres-env
file "postgres-env" created

$ cat postgres-env
DB_HOST=localhost
DB_PORT=5432
```

#### Decode directly from file

```bash
$ ls
token-secret.yaml

$ cat token-secret.yaml
apiVersion: v1
data:
  TOKEN: OGZkNDE5NzNhY2FjMDRlNWZjNzZmZGU1NDM5YzhiOTRmMWViMTIzMw==
kind: Secret
metadata:
  name: token
type: Opaque

$ k8shhh decode -i token-secret.yaml
TOKEN=8fd41973acac04e5fc76fde5439c8b94f1eb1233
```

The decoded output can also be saved into a file, which can then be used as a
configuration file for your program.

```bash
$ k8shhh decode -i token-secret.yaml -o token
file "token" created

$ cat token
TOKEN=8fd41973acac04e5fc76fde5439c8b94f1eb1233
```

#### Using kubectl with `k8shhh decode`

`k8shhh decode` also works well with [kubectl][kubectl]. Some of the examples
include the following:

```bash
$ kubectl get secret mysecret -o yaml
apiVersion: v1
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
kind: Secret
metadata:
  creationTimestamp: 2016-01-22T18:41:56Z
  name: mysecret
  namespace: default
  resourceVersion: "164619"
  selfLink: /api/v1/namespaces/default/secrets/mysecret
  uid: cfee02d6-c137-11e5-8d73-42010af00002
type: Opaque

$ kubectl get secret mysecret -o yaml | k8shhh decode
password=1f2d1e2e67df
username=admin

$ kubectl get secret mysecret -o json | k8shhh decode
password=1f2d1e2e67df
username=admin
```

#### More information

Please see [the GoDoc API page](http://godoc.org/github.com/jwangsadinata/k8shhh) for a
full API listing. For more examples, please consult `examples` directory located in each subpackages.

## Contributing

#### Bug Reports & Feature Requests

Please use the [issue tracker][issue-tracker] to report any bugs or file feature requests.

#### Developing

PRs are welcome. To begin developing, do this:

```bash
$ git clone git@github.com:jwangsadinata/k8shhh.git
$ cd k8shhh/
$ go build -mod=vendor ./cmd/k8shhh
$ ./k8shhh
usage: k8shhh [<flags>] <command> [<args> ...]

k8shhh: Quickly encode your configuration into K8s secrets.

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  encode [<flags>]
    encode your configuration as k8s secrets

  decode [<flags>]
    decode your k8s secrets into a readable format

  version
    print the current version of k8shhh.
```

This project uses [go modules][go-modules] for managing dependencies, which
comes with Go 1.11 and above. After adding a new dependency, please run the following:

```bash
$ GO111MODULE=on go mod tidy
$ GO111MODULE=on go mod vendor
```

[go-modules]: https://github.com/golang/go/wiki/Modules
[go-project]: https://golang.org/project
[issue-tracker]: https://github.com/jwangsadinata/k8shhh/issues
[kubectl]: https://kubernetes.io/docs/reference/kubectl/kubectl
[releases]: https://github.com/jwangsadinata/k8shhh/releases
[usage]: https://github.com/jwangsadinata/k8shhh#usage
