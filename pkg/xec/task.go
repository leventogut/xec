package xec

import (
	"context"
	"fmt"
	"leventogut/xec/pkg/output"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	DefaultTimeout = 600
	o              = output.GetInstance()
)

// Execute starts the defined command with it;s arguemnts.
func Execute(wg *sync.WaitGroup, taskPointerAddress **Task) {
	defer wg.Done()
	t := *taskPointerAddress
	var cancel context.CancelFunc
	if t.Timeout == 0 {
		t.Timeout = DefaultTimeout

		// o.Dev(fmt.Sprintf("Default timeout in task config not set, using global default timeout: %v", DefaultTimeout))
	}
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()

	// Merge args from config and user entered
	args := append(t.Args, t.ExtraArgs...)

	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, t.Cmd, args...)

	// Set environment values
	t.Status.ExecCmd.Env = t.SetEnvironment()

	// TODO if quiet flag is set do not log to console.
	t.Status.ExecCmd.Stdin = os.Stdin
	t.Status.ExecCmd.Stdout = os.Stdout
	t.Status.ExecCmd.Stderr = os.Stderr
	o.Info("Task " + t.Alias + " is starting")
	t.Status.Started = true

	// Execute command
	if err := t.Status.ExecCmd.Run(); err != nil {
		o.Error(fmt.Sprintf("Error: %+v\n", err))
		t.Status.Success = false
	} else {
		t.Status.Success = true
	}

	t.Status.Finished = true
	o.Info("Task " + t.Alias + " is finished")
	t.Status.ExitCode = t.Status.ExecCmd.ProcessState.ExitCode()

	// o.Dev(fmt.Sprintf("PID: %v, ExitCode: %v\n", t.Status.ExecCmd.ProcessState.Pid(), t.Status.ExecCmd.ProcessState.ExitCode()))
	if t.Status.Success {
		o.Success("Task completed successfully")
	} else {
		o.Error("Task didn't completed.")
	}
	if !t.IgnoreError && t.Status.ExitCode != 0 {
		o.Error("IgnoreError is not set and task errored. Exiting (default)")
		os.Exit(1)
	}
}

// SetEnvironment prepares the environment values using pre-defined rules based on regex in the configuration.
func (t *Task) SetEnvironment() []string {
	var environmentValuesAfterAcceptFilter []string
	var environmentValuesToBeFedToProcess []string
	var environmentValuesConfig []string

	if t.Environment.PassOn {
		// Traverse the environment values we have and apply the accept filters.
		for _, envKeyValue := range os.Environ() {
			if t.Environment.AcceptFilterRegex != nil {
				// l.Debug("AcceptFilterRegex is set")
				for _, regex := range t.Environment.AcceptFilterRegex {
					if CheckRegex(envKeyValue, regex) {
						environmentValuesAfterAcceptFilter = append(environmentValuesAfterAcceptFilter, envKeyValue)
					}
				}
			}
		}

		// Traverse the accepted environment values above and apply the reject filters.
		for _, envKeyValue := range environmentValuesAfterAcceptFilter {
			if t.Environment.RejectFilterRegex != nil {
				for _, regex := range t.Environment.RejectFilterRegex {
					if !CheckRegex(envKeyValue, regex) {
						environmentValuesToBeFedToProcess = append(environmentValuesToBeFedToProcess, envKeyValue)
					}
				}
			}
		}
	}

	// Traverse the configured environment values and add them to the list.
	for _, envKeyValueMap := range t.Environment.Values {
		envKeyValue := ConvertEnvMapToEnvString(envKeyValueMap)
		environmentValuesConfig = append(environmentValuesConfig, envKeyValue)
	}
	environmentValuesToBeFedToProcess = append(environmentValuesToBeFedToProcess, environmentValuesConfig...)

	return environmentValuesToBeFedToProcess
}
