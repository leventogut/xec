# Simple example

## Code

```yaml
tasks:
  - alias: mycommand
    cmd: echo 
    args:
        - "my command is run"
    description: run my command
```

## Description

The simplest task.

## Execution

```bash
❯ go run main.go --verbose --config examples/simple.xec.yaml mycommand
2024-02-06T19:41:46+01:00 | [INFO] | Task mycommand is starting
my command is run
2024-02-06T19:41:46+01:00 | [INFO] | Task mycommand finished in 2.443334ms.
2024-02-06T19:41:46+01:00 | [SUCCESS] | Task mycommand completed successfully in 2.443334ms.

```
