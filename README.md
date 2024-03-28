# yamlak

yamlak is a command line tool for accessing and manipulating yaml files. It consists of 3 commands that are used to get, set and delete nodes in yaml files.

## Usage
### yamlak get
Use `yamlak get` to find what value a node has:
```console
foo@bar:~$ yamlak get spec.template.spec.containers[0].image file.yaml
nginx
```
### yamlak set
You can use `yamlak set` to set a node to given value.
```console
foo@bar:~$ yamlak get spec.template.spec.containers[0].image "nginx" file.yaml
apiVersion: apps/v1
kind: Deployment
...
```
### Creating paths with yamlak set
Command `yamlak set` only works if the given node exists, otherwise it returns an error. If you want to create a path to node you should use flag `--force` or `-f`:
```console
foo@bar:~$ yamlak get spec.template.spec.containers[0].image "nginx" file.yaml --force
apiVersion: apps/v1
kind: Deployment
...
```
Flag `-f` creates the whole path to a node if it doesn't exist, but it does not work when trying to create more than one new array element. So if in spec.template.spec.containers you have just one container defined, then this command will be ok
```console
foo@bar:~$ yamlak get spec.template.spec.containers[1].image "nginx" file.yaml --force
...
```
but this one will fail
```console
foo@bar:~$ yamlak get spec.template.spec.containers[2].image "nginx" file.yaml --force
...
```

### Updating the file
By default `yamlak set` will output result of its operation to standard output. If instead you want to overwrite the current file you should use `--in-place` flag or its shorthand `-i`:
```console
foo@bar:~$ yamlak get spec.template.spec.containers[1].image "nginx" file.yaml --in-place
```

### Deleting nodes
You can delete nodes in yaml files using `yamlak delete` or `yamlak del`:
```console
foo@bar:~$ yamlak del spec.template.spec.containers[1].image file.yaml
apiVersion: apps/v1
kind: Deployment
...
```
`yamlak delete` by default prints output to standard output but it supports the `--in-force` flag when you want to overwrite the file itself

### Multi-object files
When working with mult-object (separated with "---") you can make use of ` --condition` flag or its shorthand `-c` which can help you apply your command to only some of declared objects:

```console
foo@bar:~$ yamlak del spec.template.spec.containers[0].image file.yaml
busybox
nginx
```

```console
foo@bar:~$ yamlak get spec.template.spec.containers[0].image file.yaml --condition="kind==Deployment"
nginx
```

```console
foo@bar:~$ yamlak get spec.template.spec.containers[0].image file.yaml --condition="spec.replicas>1"
nginx
```
The condition flag supports operators ">", "<", "<=", ">=", "!=" oraz "==".
When the `--condition` flag is used with `yamlak set` or `yamlak del` it will only modify objects that fulfill condition but the output will consist of all of objects in file.

### Using docker image
You can use yamlak by using its docker image but you have to remember to mount the file you want to parse as docker volume:
```console
foo@bar:~$ docker run -v file.yaml:/var/file.yaml slimo300/yamlak get spec.template.spec.containers[0].image /var/file.yaml
```