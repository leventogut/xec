# Restart on failure

## Code

```yaml
tasks:
  - alias: lsz
    cmd: ls
    args:
      - "-z"
    restartOnFailure: true
```

## Description

## Execution

```bash
❯ xec --verbose --config examples/restart-on-error.xec.yaml lsz
2024-02-03T22:06:16+01:00 | [INFO] | Task lsz is starting
ls: invalid option -- z
usage: ls [-@ABCFGHILOPRSTUWabcdefghiklmnopqrstuvwxy1%,] [--color=when] [-D format] [file ...]
2024-02-03T22:06:16+01:00 | [ERROR] | Error: exit status 1

2024-02-03T22:06:16+01:00 | [INFO] | Task lsz finished in 2.442792ms.
2024-02-03T22:06:16+01:00 | [ERROR] | Task lsz didn't completed.
2024-02-03T22:06:16+01:00 | [INFO] | Task lsz is starting
ls: invalid option -- z
usage: ls [-@ABCFGHILOPRSTUWabcdefghiklmnopqrstuvwxy1%,] [--color=when] [-D format] [file ...]
2024-02-03T22:06:16+01:00 | [ERROR] | Error: exit status 1

2024-02-03T22:06:16+01:00 | [INFO] | Task lsz finished in 4.519292ms.
2024-02-03T22:06:16+01:00 | [ERROR] | Task lsz didn't completed.
2024-02-03T22:06:16+01:00 | [INFO] | Task lsz is starting
ls: invalid option -- z
usage: ls [-@ABCFGHILOPRSTUWabcdefghiklmnopqrstuvwxy1%,] [--color=when] [-D format] [file ...]
2024-02-03T22:06:16+01:00 | [ERROR] | Error: exit status 1

2024-02-03T22:06:16+01:00 | [INFO] | Task lsz finished in 4.671417ms.
2024-02-03T22:06:16+01:00 | [ERROR] | Task lsz didn't completed.
...
...
...
```
