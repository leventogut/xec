# Parallel execution

## Code

```yaml
tasks:
  - alias: wait_10
    cmd: sleep
    args:
      - "10"
  - alias: wait_3
    cmd: sleep
    args:
      - "3"
taskLists:
  - alias: parallel
    parallel: true
    taskAliases:
      - wait_10
      - wait_3
```

## Description

Here we have two tasks defined; wait_10 and wait_3. Also there is a task list called parallel. Parallel is configured to run it's tasks in parallel using `parallel: true` stanza.

## Execution

TaskList timestamp shows that the tasks in the lists have finished in ~10s.

```bash
❯ xec --verbose --config examples/parallel.xec.yaml parallel
2024-02-03T21:49:59+01:00 | [INFO] | Task wait_3 is starting
2024-02-03T21:49:59+01:00 | [INFO] | Task wait_10 is starting
2024-02-03T21:50:02+01:00 | [INFO] | Task wait_3 finished in 3.009180125s.
2024-02-03T21:50:02+01:00 | [SUCCESS] | Task wait_3 completed successfully in 3.009180125s.
2024-02-03T21:50:09+01:00 | [INFO] | Task wait_10 finished in 10.008260667s.
2024-02-03T21:50:09+01:00 | [SUCCESS] | Task wait_10 completed successfully in 10.008260667s.
2024-02-03T21:50:09+01:00 | [SUCCESS] | TaskList parallel finished in 10.00922625s.
```
