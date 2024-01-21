<!-- vale Microsoft.HeadingAcronyms = NO -->
# Xec

Xec is a simple command (task) executor, that you define your tasks as commands and arguments in a yaml file, and it executes them when you call xec with the alias you defined.

It allows you to control environment that is passed on to the sub-process, such as timeout, environment values and so on.

Environment values flow can be defined, to pass or block values based on regex, also adding environment values is trivial.
Xec supports reading .env file.

Usage:

To see all available aliases just enter with no arguments.

```bash
xec
```

To run an alias, just enter name of the alias and additional parameters if required. Arguments after "--" is appended to the tasks' config arguments.

```bash
xec myls
```

Or a more advanced usage:

```bash
xec --ignore-errors myls -- -h
```

## Defaults

- IgnoreErrors: false
- Timeout: Timeout for a task. -> 10 minutes.

## Parallelism

TaskList has the option `parallel`, when set to true xec will run the tasks in parallel.

```yaml
  - alias: parallel
    parallel: true
    taskNames:
      - wait_10
      - wait_5
```

## A task

All configuration options of a task:

```yaml
tasks:
  - alias: myls
    description: Execute ls with params.
    cmd: ls
    args:
      - "-al"
      - "-h"
    timeout: 10 # Timeout is 10 seconds
    environment:
      passOn: true # Pass the environment key/values that Xec receives to the process or not.
      values: # Additional environment values to pass to the process.
        key: value
        environment: prod
    acceptFilterRegex: # Include the environment values matches the regex.
      - "XEC_*"
    rejectFilterRegex: # Don't include the environment values matches the regex.
      - "SECRET*"
      - "AWS*"
    debug: true # Enable debug on Xec about this task.
    logFile: "myls.log" # Log to this file.
    ignoreError: true # Ignore if task is errored.
```

## Error handling of tasks

You can ignore (and continue) errored tasks. This can be achieved in three levels, TaskDefaults, Task, TaskLists.
TaskDefaults set's all task instances, while Task level affects individual. The other level affects TaskList which affects all tasks in the task list.

## Writing configuration files (schema)

JSON schema can be found [here](https://test.test)

## Examples

### Simple example

```yaml
tasks:
  - alias: myls
    description: execute custom ls.
    cmd: ls
    args:
      - -al
```

### Task list

```yaml
...
taskLists:
  - alias: lsenv
    description: "tasklist for ls and env"
    taskNames:
      - ls
      - printenv
  - alias: lszenv
    description: "tasklist for ls and env errors"
    taskNames:
      - lsz
      - printenv
    ignoreError: true
```

## Contributing

Contributions are most welcome.

## Build

## Install
