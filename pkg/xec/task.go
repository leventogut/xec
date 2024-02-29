package xec

import (
	"context"
	"github.com/leventogut/xec/pkg/output"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
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

// Execute starts the defined command with its arguments.
func Execute(taskPointerAddress **Task) {
	t := *taskPointerAddress

	var cancel context.CancelFunc
	if t.Timeout == 0 {
		t.Timeout = DefaultTimeout
	}
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()

	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, t.Cmd, append(t.Args, t.ExtraArgs...)...)

	// Change working directory if it is given.
	if t.Directory != "" {
		t.Status.ExecCmd.Dir = t.Directory
	}

	// Set environment values
	t.Status.ExecCmd.Env = t.SetEnvironment()

	t.Status.ExecCmd.Stdin = os.Stdin

	// Setup log
	var logFile *os.File

	if t.LogFile != "" {
		var err error

		logFile, err = os.OpenFile(t.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			o.Error("Can't open log file %+v - Error: %+v\n", t.LogFile, err)
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

	o.Info("Task %+v is starting.", t.Alias)
	o.Info("Task %+v is logged to %+v", t.Alias, t.LogFile)
	t.Status.Started = true
	taskStartTime := time.Now()

	// Start execution
	if err := t.Status.ExecCmd.Start(); err != nil {
		o.Error("Task couldn't be started, Error: %+v\n", err)
		os.Exit(1)
	}

	// Handle signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for {
			receivedSignal := <-signalChannel
			o.Warning("Signal received: %+v.", receivedSignal)
			o.Info("Passing signal %+v to %+v.", receivedSignal, t.Alias)
			_ = t.Status.ExecCmd.Process.Signal(receivedSignal)
		}
	}()

	//Wait for the execution
	if err := t.Status.ExecCmd.Wait(); err != nil {
		o.Error("Error occurred while waiting, Error: %+v\n", err)
		os.Exit(1)
	}

	t.Status.Finished = true
	taskFinishTime := time.Now()
	taskDuration := taskFinishTime.Sub(taskStartTime)

	o.Info("Task " + t.Alias + " finished in " + taskDuration.String() + ".")

	if t.Status.ExecCmd.ProcessState.ExitCode() > 0 {
		t.Status.Success = false
		o.Error("Task " + t.Alias + " didn't complete successfully.")

	} else if t.Status.ExecCmd.ProcessState.ExitCode() == 0 {
		t.Status.Success = true
		o.Success("Task " + t.Alias + " completed successfully in " + taskDuration.String() + ".")
	}

	// Restarts
	if t.RestartOnSuccess && t.Status.Success {
		t.NumberOfRestarts++
		if t.NumberOfRestarts < t.RestartLimit {
			Execute(&t)
		}

	}
	if t.RestartOnFailure && !t.Status.Success {
		t.NumberOfRestarts++
		if t.NumberOfRestarts < t.RestartLimit {
			Execute(&t)
		}
	}

	if !t.IgnoreError && t.Status.ExecCmd.ProcessState.ExitCode() != 0 {
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
