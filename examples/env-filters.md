# Environment filters

## Code

```yaml
verbose: true
tasks:
  - alias: penv
    cmd: printenv
    environment:
      values:
        - XEC_task: penv
        - XEC_environment: dev
      passOn: true
      acceptFilterRegex:
        - "XEC.*"
      rejectFilterRegex:
        - ".*SECRET.*"
```

## Description

## Execution

Setup

Export some environment key/values for testing.

```bash
export XEC_test=1 
export XEC_SECRET_1=mySecret
```
### Before application of reject filter (accept filter is active)

Accept filter accepts only environment keys that start with XEC. The two environment key/values that start with xec_ is coming from the configuration.

```bash
❯ xec --config examples/env-filters.xec.yaml penv                                                            ☸ docker-desktop in xec on  shell-support [!+] via 🐹 v1.22.0 with unknown env 
2024-03-08T00:32:20+01:00 | SUCCESS | Loaded config file, [examples/env-filters.xec.yaml]
2024-03-08T00:32:20+01:00 | SUCCESS | Loaded config file, [$HOME/.xec.yaml]
2024-03-08T00:32:20+01:00 | INFO | Task penv is starting.
2024-03-08T00:32:20+01:00 | INFO | Task penv is not logged.
XEC_test=1
XEC_SECRET_1=mySecret
xec_task=penv
xec_environment=dev
2024-03-08T00:32:20+01:00 | INFO | Task penv finished in 1.789458ms.
2024-03-08T00:32:20+01:00 | SUCCESS | Task penv completed successfully in 1.789458ms.
```

### After enabling the reject filter

This time XEC_SECRET environment key/value is not present on the output.
```bash
❯ xec --config examples/env-filters.xec.yaml penv                                                            ☸ docker-desktop in xec on  shell-support [!+] via 🐹 v1.22.0 with unknown env 
2024-03-08T00:36:49+01:00 | SUCCESS | Loaded config file, [examples/env-filters.xec.yaml]
2024-03-08T00:36:49+01:00 | SUCCESS | Loaded config file, [$HOME/.xec.yaml]
2024-03-08T00:36:49+01:00 | SUCCESS | Loaded config file, [./.xec.yaml]
2024-03-08T00:36:49+01:00 | INFO | Task penv is starting.
2024-03-08T00:36:49+01:00 | INFO | Task penv is not logged.
XEC_test=1
xec_task=penv
xec_environment=dev
2024-03-08T00:36:49+01:00 | INFO | Task penv finished in 1.964958ms.
2024-03-08T00:36:49+01:00 | SUCCESS | Task penv completed successfully in 1.964958ms.
```