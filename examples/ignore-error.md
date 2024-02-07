# Title of the example

## Code

```yaml
tasks:
  - alias: printenv
    description: execute printenv
    cmd: printenv
    environment:
      values:
        - task: printenv
        - env: stg
  - alias: lsz
    description: ls to give error.
    cmd: ls
    args:
      - "-z"
taskLists:
  - alias: lszenv
    taskAliases:
      - lsz
      - printenv
    ignoreError: false
```

## Description

## Execution

### ignoreError is set to false (default)

```bash
❯ xec --verbose --config examples/ignore-error.xec.yaml lszenv
2024-02-06T19:54:37+01:00 | [INFO] | Task lsz is starting
ls: invalid option -- z
usage: ls [-@ABCFGHILOPRSTUWabcdefghiklmnopqrstuvwxy1%,] [--color=when] [-D format] [file ...]
2024-02-06T19:54:37+01:00 | [ERROR] | Error: exit status 1

2024-02-06T19:54:37+01:00 | [INFO] | Task lsz finished in 3.288292ms.
2024-02-06T19:54:37+01:00 | [ERROR] | Task lsz didn't completed.
exit status 1
```

### ignoreError is set to true

```bash
❯ xec --verbose --config examples/ignore-error.xec.yaml lszenv
2024-02-06T19:53:25+01:00 | [INFO] | Task lsz is starting
ls: invalid option -- z
usage: ls [-@ABCFGHILOPRSTUWabcdefghiklmnopqrstuvwxy1%,] [--color=when] [-D format] [file ...]
2024-02-06T19:53:25+01:00 | [ERROR] | Error: exit status 1

2024-02-06T19:53:25+01:00 | [INFO] | Task lsz finished in 2.561542ms.
2024-02-06T19:53:25+01:00 | [ERROR] | Task lsz didn't completed.
2024-02-06T19:53:25+01:00 | [INFO] | Task printenv is starting
task=printenv
environment2=dev2
2024-02-06T19:53:25+01:00 | [INFO] | Task printenv finished in 1.745458ms.
2024-02-06T19:53:25+01:00 | [SUCCESS] | Task printenv completed successfully in 1.745458ms.
2024-02-06T19:53:25+01:00 | [SUCCESS] | TaskList lszenv finished in 4.965958ms.
```
