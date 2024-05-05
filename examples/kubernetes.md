# Kubernetes example

## Code

```yaml
namespace: k
tasks:
  - alias: gp
    cmd: kubectl
    args:
      - get
      - pods
  - alias: gs
    cmd: kubectl
    args:
      - get
      - svc
taskLists:
  - alias: ps
    taskAliases:
      - gp
      - gs

```

## Description

`namespace` is set to `k`.

## Execution

```bash
❯ xec --verbose --config examples/kubernetes.xec.yaml k ps

```
