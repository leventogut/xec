<!-- vale Microsoft.HeadingAcronyms = NO -->
# Xec

Xec is a simple task executor, that you define your tasks as commands and arguments, and it executes them when you call xec with the alias you defined.
It allows you control some elements of the environment that your tasks run, such as timeout, environment values and so on.

Usage:

```bash
xec --debug --verbose --dev ls -- -h
```

## TODO List

- [X] Flags are not working (current solution dictates PFlags are not used by the sub-processes, not acceptable)
- [ ] Multiple config file reading.
- [X] Find which license to use and copy it.
- [ ] Start writing tests.
- [ ] Importing function.
- [ ] reading .env
- [ ] Ability to add extra args in TaskLists?
- [X] --quiet is not working

## Features

- [?] Reads .env by default.
- [ ] Can import multiple config/task files.
- [X] Environment value security with pass on/off, accept and, reject filters.
- [X] Control environment value flow, including regex based accept and reject filters.
- Define timeout for task
- Easily create new aliases for lists of tasks
- Ability of passing extra arguments to the task itself.
- Label and tag tasks to run a bunch of tasks based on.
- Ability import other Xec files as new aliases.
- Separate logging for each task
- Setting default values of options in TaskDefaults
- Ability add environment values in defaults and individual task.
- Looks for a file named .xec.yaml in either current directory or $HOME.

## Defaults

- Dev: Development mode. -> false
- Debug: Debug outputs.  -> false
- Verbose: Additional informational outputs. -> false
- Quiet: Disable any console output. -> false
- Timeout: Timeout for a task. -> 10 seconds.

## Behaviors

- Environment values defined in TaskDefault and individual task are merged.

## Full configuration options of a task

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
      passOn: true # Pass the environment key/values Xec receives to the process or not.
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
```

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
