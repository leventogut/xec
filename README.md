<!-- vale Microsoft.HeadingAcronyms = NO -->
# xec

xec is a simple command executor.

Command, it's arguments and configuration of how to run it, is referred as `task`.
xec reads a yaml based configuration file in either current directory or in home directory of the user.
Reading configuration from the current working directory has advantages of per-project/per-task configuration structure.

It allows you to control environment that is passed on to the command, such as timeout, environment values, restart behavior and so on.

Environment values flow can be defined, to pass or block values based on regex, also adding environment values is trivial. Also xec supports reading .env file.

xec has the capability of:

- Adding extra arguments via cli
- Grouping tasks as task lists
- Run tasks in parallel (via task lists)
- Restart task based on exit code, failure or success
- Filtering and adding environment values that are passed to the command

## Table of contents

- [xec](#xec)
  - [Table of contents](#table-of-contents)
  - [Usage](#usage)
  - [Defaults](#defaults)
    - [Parallelism](#parallelism)
  - [Anatomy of a task](#anatomy-of-a-task)
  - [Error handling of tasks](#error-handling-of-tasks)
  - [Writing configuration files (schema)](#writing-configuration-files-schema)
  - [Examples](#examples)
  - [Contributing](#contributing)
  - [Build](#build)
    - [Release build](#release-build)
    - [Snapshot build](#snapshot-build)
  - [Install](#install)

## Usage

To see all available aliases just enter with no alias argument.

```bash
❯ xec
Simple command executor.

Usage:
  xec <flags> <alias> -- [args-to-be-passed] [flags]
  xec [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  lsz         ls-z
  lszenv      
  printenv    printenv
  version     Print the version number

Flags:
      --config string     config file to read (default is ~/.xec.yaml,  $PWD/.xec.yaml)
      --debug             Debug level output.
  -h, --help              help for xec
      --ignore-error      Ignore errors on tasks.
      --log-file string   Filename to use for logging.
      --no-color          Disable color output.
      --quiet             No output except errors].
      --verbose           Verbose level output.

Use "xec [command] --help" for more information about a command.
exit status 1
```

To run a task, just enter name of the alias (and additional parameters if required).

Arguments after "--" is appended to the tasks' config arguments.

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

### Parallelism

TaskList has the option `parallel`, when set to true xec will run the tasks in parallel.

```yaml
taskLists:
  - alias: parallel
    parallel: true
    taskNames:
      - wait_10
      - wait_5
```

## Anatomy of a task

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

JSON schema can be found [here](https://raw.githubusercontent.com/leventogut/xec/main/schema/xec-tasks-yaml-schema.json)

## Examples

| ------------------ | ---------------------------------------------- | ------------------------------------------- |
| Parallel execution | [documentation](examples/parallel.md)          | [code](examples/parallel.xec.yaml)          |
| Restart on failure | [documentation](examples/restart-on-failure.md)| [code](examples/restart-on-failure.xec.yaml)|
| Restart on success | [documentation](examples/restart-on-success.md)| [code](examples/restart-on-success.xec.yaml)|

## Contributing

Contributions are most welcome.

## Build

### Release build

Release build is done with `goreleaser` in GH Actions.

### Snapshot build

```bash
goreleaser release --snapshot --clean
```

OR

```bash
xec build
```

## Install
