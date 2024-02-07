package xec

import (
	"context"
	"fmt"
	"io"
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

func ExecuteWithWaitGroups(wg *sync.WaitGroup, taskPointerAddress **Task) {
	defer wg.Done()
	Execute(taskPointerAddress)
}

// Execute starts the defined command with it's arguments.
func Execute(taskPointerAddress **Task) {
	t := *taskPointerAddress
	var cancel context.CancelFunc
	if t.Timeout == 0 {
		t.Timeout = DefaultTimeout
	}
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()

	// Merge args from config and user entered
	args := append(t.Args, t.ExtraArgs...)

	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, t.Cmd, args...)

	// Change working directory if it is given.
	if t.Directory != "" {
		t.Status.ExecCmd.Dir = t.Directory
	}

	// Set environment values
	t.Status.ExecCmd.Env = t.SetEnvironment()

	t.Status.ExecCmd.Stdin = os.Stdin
	var logFile *os.File

	if t.LogFile != "" {
		var err error

		logFile, err = os.OpenFile(t.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			o.Error(fmt.Sprintf("Can't open log file %+v - Error: %+v\n", t.LogFile, err))
		}
		defer logFile.Close()

	} else {
		logFile, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	}

	if o.QuietFlag() {
		t.Status.ExecCmd.Stdout = io.MultiWriter(logFile)
		t.Status.ExecCmd.Stderr = io.MultiWriter(logFile)

	} else {
		t.Status.ExecCmd.Stdout = io.MultiWriter(logFile, os.Stdout)
		t.Status.ExecCmd.Stderr = io.MultiWriter(logFile, os.Stderr)

	}

	o.Info("Task " + t.Alias + " is starting")
	t.Status.Started = true
	taskStartTime := time.Now()

	// Execute command
	if err := t.Status.ExecCmd.Run(); err != nil {
		o.Error(fmt.Sprintf("Error: %+v\n", err))
		t.Status.Success = false
	} else {
		t.Status.Success = true
	}

	t.Status.Finished = true
	taskFinishTime := time.Now()
	taskDuration := taskFinishTime.Sub(taskStartTime)

	o.Info("Task " + t.Alias + " finished in " + taskDuration.String() + ".")
	t.Status.ExitCode = t.Status.ExecCmd.ProcessState.ExitCode()

	if t.Status.Success {
		o.Success("Task " + t.Alias + " completed successfully in " + taskDuration.String() + ".")
	} else {
		o.Error("Task " + t.Alias + " didn't completed.")
	}
	if t.RestartOnSuccess && t.Status.Success {
		Execute(&t)
	}

	if t.RestartOnFailure && !t.Status.Success {
		Execute(&t)
	}
	if !t.IgnoreError && t.Status.ExitCode != 0 {
		os.Exit(t.Status.ExitCode)
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
