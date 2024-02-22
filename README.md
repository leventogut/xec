<!-- vale Microsoft.HeadingAcronyms = NO -->
# xec a simple command executor

xec is a simple command executor.

Command, it's arguments and configuration of how to run it, is referred as `task`.
Xec reads a yaml based configuration file in either current directory and in home directory of the user. It is also possible to add another config file via `--config` argument.
Reading configuration from the current working directory has advantages of per-project configuration structure.

It allows you to control environment of the command execution, such as timeout, environment values, restart behavior and so on.

Environment values can be filtered, either to pass or block values based on regex match, also adding environment values is trivial. 
Xec also supports reading .env file in the current directory, and reads it by default.

xec has the capability of:

- Dynamic sub-command generation based on the aliases of tasks
- Adding extra arguments via cli
- Grouping tasks as task lists  
- Run tasks in parallel (via task lists)
- Restart task based on exit code, failure or success
- Filtering and adding environment values that are passed to the command
- Importing multiple configurations


- [xec a simple command executor](#xec-a-simple-command-executor)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Defaults](#defaults)
  - [Dependencies](#dependencies)
  - [Initial configuration](#initial-configuration)
  - [Examples](#examples)
  - [Anatomy of a task](#anatomy-of-a-task)
  - [Error handling of tasks](#error-handling-of-tasks)
  - [Writing configuration files (schema)](#writing-configuration-files-schema)
  - [Contributing](#contributing)
  - [Build](#build)
    - [Release build](#release-build)
    - [Snapshot build](#snapshot-build)
    - [Releasing xec](#releasing-xec)

## Installation

### Pre-built binaries

Visit the [releases page](https://github.com/leventogut/xec/releases) or [latest release page](https://github.com/leventogut/xec/releases/latest).

#### Linux and Darwin

```bash
# set the version
export XEC_VERSION=v0.0.10
# Download the archive suits to your OS and architecture.
curl -L -o xec.tar.gz https://github.com/leventogut/xec/releases/download/${XEC_VERSION}/xec_Linux_x86_64.tar.gz
# Extract the archive
tar -xzvf xec.tar.gz
# Move the binary (might need root/sudo, depending on the path)
mv xec /usr/local/bin
```

#### Verification

```bash
# Set the version
export XEC_VERSION=v0.0.10
# Download checksums file
curl -L -o xec_checksums.txt https://github.com/leventogut/xec/releases/download/${XEC_VERSION}/xec_${XEC_VERSION}_checksums.txt
# Get checksum
sha256sum xec.tar.gz
# Compare the checksum from the file
```


### Go

```bash
go install github.com/leventogut/xec 
```

## Usage

To see all available aliases just enter with no argument, it will read all configurations and generate the sub-commands/aliases.
In the following output `build` and `env` are aliases available.

Description of an alias command constitutes the command to be run and its arguments.

```bash
❯ xec                                                                                                                                                                                   xec on  main [!?] via 🐳 desktop-linux via 🐹 v1.20.4 with unknown env 
Simple command executor.

Usage:
  xec <flags> <alias> -- [additional-args] [flags]
  xec [command]

Available Commands:
  build       goreleaser release --snapshot --clean
  completion  Generate the autocompletion script for the specified shell
  env         printenv 
  help        Help about any command
  init        initialize a configuration file in the current directory.
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

If no arguments are given, it exits with an exit code of 1.

To run a task, just enter name of the alias (and additional parameters if required).

```bash
xec myls
```

Arguments after "--" is appended to the tasks' config arguments.

An example with additional arguments:

```bash
xec --ignore-errors myls -- -h
```

## Defaults

- Verbose: true
- Debug: false
- Config: ""
- NoColor: false
- Quiet: false
- IgnoreErrors: false
- Timeout: 600s

## Dependencies

xec uses [Cobra](https://github.com/spf13/cobra) for it's command generation.

## Initial configuration

An initial, skeleton configuration can be created via:

```bash
❯ xec init                                                                                                                                                                              xec on  main [✘!?] via 🐳 desktop-linux via 🐹 v1.20.4 with unknown env 
2024-02-10T21:33:54+01:00 | [SUCCESS] | Init configuration is written to file .xec.yaml.
```
# Examples

| Feature               | Documentation                                                    | Code                                                                         |
|-----------------------|------------------------------------------------------------------|------------------------------------------------------------------------------|
| Simple                | [examples/simple.md](examples/parallel.md)                       | [examples/simple.xec.yaml](examples/parallel.xec.yaml)                       |
| Parallel execution    | [examples/parallel.md](examples/parallel.md)                     | [examples/parallel.xec.yaml](examples/parallel.xec.yaml)                     |
| Ignore error          | [examples/ignore-error.md](examples/parallel.md)                 | [examples/ignore-error.xec.yaml](examples/parallel.xec.yaml)                 |
| Restart on failure    | [examples/restart-on-failure.md](examples/restart-on-failure.md) | [examples/restart-on-failure.xec.yaml](examples/restart-on-failure.xec.yaml) |
| Restart on success    | [examples/restart-on-success.md](examples/restart-on-success.md) | [examples/restart-on-success.xec.yaml](examples/restart-on-success.xec.yaml) |
| Import Configurations | [examples/import.md](examples/restart-on-success.md)             | [examples/import.xec.yaml](examples/restart-on-success.xec.yaml)             |

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
    restartOnSuccess: false
    restartOnFailure: false
```

## Error handling of tasks

You can ignore (and continue) an errored tasks. This can be achieved in three levels, TaskDefaults, Task, TaskLists.
TaskDefaults set's all task instances, while Task level affects individual. The other level affects TaskList which affects all tasks in the task list.

Restarting the task based on the exit code is supported as well, with restartOnSuccess and restartOnFailure

## Writing configuration files (schema)

JSON schema can be found [here](https://raw.githubusercontent.com/leventogut/xec/main/schema/xec-tasks-yaml-schema.json)

## Contributing

Contributions are most welcome.
Please create a feature branch and work on that. Once your feature is ready raise PR.

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

### Releasing xec

When there is a version tag attached, build and release is automatically done.

```bash
git commit -m "doing some stuff related to ..."
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```
