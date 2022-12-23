<!-- vale Microsoft.HeadingAcronyms = NO -->
# Xec

## Features

- Reads .env by default.
- Can import other `.xec.yaml` files

Looks for ./.xec.(yaml|yml) if not found it traverses up to /

## Base example

```yaml
build:

```

## Multiline example

```yaml
docker-build-release:
  - check-docs
  - |+  
  docker build \
      -f deploy/skaffold/Dockerfile.lts \
      --target release \
      -t gcr.io/$(GCP_PROJECT)/skaffold:edge-lts \
      -t gcr.io/$(GCP_PROJECT)/skaffold:$(COMMIT)-lts \
      .
```

## Example project YAML

```yaml
# examples/.xec.yaml
myTask:
cmd: "ls -alH"
read:
env:
- key: filename
    value: README.md
pre-cmd: ls -al README.md
cmd: cat ${filename}
```

```yaml
# ~/dev/org/simple/.xec.yaml
tasks:
  one: echo "one"
  two: echo "two"
  three: echo "three"
```

```shell
xec add task myTask -- "ls -alH"
```

```shell
xec MyTask
```


## Example root YAML

```yaml

```
