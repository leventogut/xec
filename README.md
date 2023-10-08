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

## TODO List

- [ ] reading .env
- [ ] In environment values section, if no filter is defined no external environment value is received.
- [ ] env-file (.env) flag and check if we can change it in the library we are using.

## Features

- [ ] Reads .env by default.
- [ ] Can import multiple xec files under self-defined alias.
- [X] Control environment value security with pass on/off, accept and, reject filters.
- [X] Define timeout for a task.
- [X] Easily create new aliases for lists of tasks that are previously defined.
- [X] Setting default values of task options with TaskDefaults section in configuration file.
- [X] Looks for a file named .xec.yaml in either current and $HOME directory.

## Planned features

- [ ] Different logging for each task/task list as configurable i.e. logFile field for each task and task list.
- [ ] Importing multiple config files function.
- [ ] Ability to add extra args in TaskLists?

## Defaults

- IgnoreErrors: false
- Timeout: Timeout for a task. -> 10 seconds.

## Behaviors

- Environment values defined in TaskDefault and individual task are merged.

## All fields of a task

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
TaskDefaults set's all task instances, while Task level affects individual. The other level is TaskList which affects all tasks in the task list.

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

### Importing

You can import other xec files under a alias/context name.

```yaml
tasks:
  - alias: k
    description: Kubernetes tasks.
    import: ~/.xec/skills/xec-kubernetes.yaml
```

## Contributing

Contributions are most welcome.
