# Restart on error

## Code

```yaml
tasks:
  - alias: myecho
    cmd: echo
    args:
      - 5
    restartOnSuccess: true

```

## Description

## Execution

```bash
❯ go run main.go --verbose --config examples/restart-on-success.xec.yaml
2024-02-06T18:58:55+01:00 | [INFO] | Task myecho is starting
5
2024-02-06T18:58:55+01:00 | [INFO] | Task myecho finished in 2.561ms.
2024-02-06T18:58:55+01:00 | [SUCCESS] | Task myecho completed successfully in 2.561ms.
2024-02-06T18:58:55+01:00 | [INFO] | Task myecho is starting
5
2024-02-06T18:58:55+01:00 | [INFO] | Task myecho finished in 1.257542ms.
2024-02-06T18:58:55+01:00 | [SUCCESS] | Task myecho completed successfully in 1.257542ms.
2024-02-06T18:58:55+01:00 | [INFO] | Task myecho is starting
5
...
...
...
```
