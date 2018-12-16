## Examples for k8shhh

This document outlines how one would use `k8shhh` for encoding and decoding
one's secrets. Two files are provided, namely `example-file` and
`token-secret.yaml`, which will be used as the starting points for these
examples.

### Encoding

For encoding your configuration, simply use `k8shhh encode` along with the
name of the input file, by using the `-i` flag. It should look like this:

```bash
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

If you wish to change the format, you can use the `-f` flag to specify the
desired format. The supported formats are `json` and `yaml`.

```bash
$ cat example-file
DB_HOST=localhost
DB_PORT=5432

$ k8shhh encode -i example-file -f json
{
        "apiVersion": "v1",
        "data": {
                "DB_HOST": "bG9jYWxob3N0",
                "DB_PORT": "NTQzMg=="
        },
        "kind": "Secret",
        "metadata": {
                "name": "mysecret"
        },
        "type": "Opaque"
}
```

If you wish to change the name of the secret generated, there is a helpful `-n`
flag that allows you to specify the generated secret name, as follows:

```bash
$ cat example-file
DB_HOST=localhost
DB_PORT=5432

$ k8shhh encode -i example-file -f json -n my-database
{
        "apiVersion": "v1",
        "data": {
                "DB_HOST": "bG9jYWxob3N0",
                "DB_PORT": "NTQzMg=="
        },
        "kind": "Secret",
        "metadata": {
                "name": "my-database"
        },
        "type": "Opaque"
}
```

You can also save the encoded configuration to a file. To do this, use the `-o`
flag and specify the name of the output file. The name of the file will also be
the name of the secret, if `-n` flag is not specified.

```bash
$ cat example-file
DB_HOST=localhost
DB_PORT=5432

$ k8shhh encode -i example-file -o my-database
my-database.yaml

$ cat postgres.yaml
apiVersion: v1
data:
  DB_HOST: bG9jYWxob3N0
  DB_PORT: NTQzMg==
kind: Secret
metadata:
  name: my-database
type: Opaque
```

For more information, you can always check the help page for `k8shhh encode`,
by typing the following:

```bash
$ k8shhh encode --help
usage: k8shhh encode [<flags>]

encode your configuration as k8s secrets

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -n, --name=NAME      the name of the generated secret
  -i, --input=INPUT    the name of the input file to encode (if input is not provided via STDIN)
  -o, --output=OUTPUT  the name of the file to write the output to (outputs to STDOUT by default). file extension will be automatically generated based on the format.
  -f, --format="yaml"  format of the generated secret (json or yaml, defaults to yaml)
```

### Decoding

To learn about `k8shhh decode`, the file `token-secret.yaml` has been provided.
Use the `-i` flag with the name of the file, to decode the given file:

```bash
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

You can also the decoded secret into a file, by using the `-o` flag and
specifying the name of the output file, like shown below:

```bash
$ cat token-secret.yaml
apiVersion: v1
data:
  TOKEN: OGZkNDE5NzNhY2FjMDRlNWZjNzZmZGU1NDM5YzhiOTRmMWViMTIzMw==
kind: Secret
metadata:
  name: token
type: Opaque

$ k8shhh decode -i token-secret.yaml -o token
file "token" created

$ cat token
TOKEN=8fd41973acac04e5fc76fde5439c8b94f1eb1233
```

Similarly, the help page for `k8shhh decode` should be accessible to you by the
following command:

```bash
$ k8shhh decode --help
usage: k8shhh decode [<flags>]

decode your k8s secrets into a readable format

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -i, --input=INPUT    the name of the input file to decode (if input is not provided via STDIN)
  -o, --output=OUTPUT  the name of the file to write the output to (outputs to STDOUT by default)
```

### Questions

If you have any questions or issues with these examples, please use the
[issue tracker][issue-tracker] to report your concerns. Contributions are very
much welcomed.
