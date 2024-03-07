# Task lists

## Configuration

```yaml
verbose: true
tasks:
  - alias: sleep_2
    cmd: sleep
    args:
      - "2"
  - alias: sleep_3
    cmd: sleep
    args:
      - "3"
taskLists:
  - alias: tl
    taskAliases:
      - sleep_2
      - sleep_3

```

## Description

## Execution

```bash
❯ go run main.go --config examples/tasklists.xec.yaml tl                                                                ☸ docker-desktop in xec on  shell-support [!+] via 🐹 v1.22.0 with unknown env 
2024-03-08T00:53:23+01:00 | SUCCESS | Loaded config file, [examples/tasklists.xec.yaml]
2024-03-08T00:53:23+01:00 | SUCCESS | Loaded config file, [$HOME/.xec.yaml]
2024-03-08T00:53:23+01:00 | SUCCESS | Loaded config file, [./.xec.yaml]
2024-03-08T00:53:23+01:00 | INFO | Task list tl is starting.
2024-03-08T00:53:23+01:00 | INFO | Task list tl is logged to 
2024-03-08T00:53:23+01:00 | INFO | Task sleep_2 is starting.
2024-03-08T00:53:23+01:00 | INFO | Task sleep_2 is not logged.
2024-03-08T00:53:25+01:00 | INFO | Task sleep_2 finished in 2.007440917s.
2024-03-08T00:53:25+01:00 | SUCCESS | Task sleep_2 completed successfully in 2.007440917s.
2024-03-08T00:53:25+01:00 | INFO | Task sleep_3 is starting.
2024-03-08T00:53:25+01:00 | INFO | Task sleep_3 is not logged.
2024-03-08T00:53:28+01:00 | INFO | Task sleep_3 finished in 3.006981208s.
2024-03-08T00:53:28+01:00 | SUCCESS | Task sleep_3 completed successfully in 3.006981208s.
2024-03-08T00:53:28+01:00 | SUCCESS | Task list tl finished in 5.016506708s.
```