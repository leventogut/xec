package xec

import (
	"context"
	"fmt"
	"leventogut/xec/pkg/output"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	DefaultTimeout = 600
	l              = output.Log
)

// Execute starts the defined command with it;s arguemnts.
func (t *Task) Execute() {
	fmt.Println("In Execute")
	// fmt.Printf("address in main: %p\n", t)
	fmt.Printf("%+v", t)
	fmt.Printf("Executing: %s\n", t.Cmd)
	var cancel context.CancelFunc
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, t.Cmd, t.Args...)
	t.SetEnvironment()

	t.Status.ExecCmd.Stdin = os.Stdin
	t.Status.ExecCmd.Stdout = os.Stdout
	t.Status.ExecCmd.Stderr = os.Stderr
	if err := t.Status.ExecCmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("t.Status.ExecCmd.ProcessState.ExitCode(): %v\n", t.Status.ExecCmd.ProcessState.ExitCode())
	fmt.Printf("t.Status.ExecCmd.ProcessState.Pid(): %v\n", t.Status.ExecCmd.ProcessState.Pid())
	fmt.Printf("t.Status.ExecCmd.ProcessState.SysUsage(): %v\n", t.Status.ExecCmd.ProcessState.SysUsage())
	fmt.Printf("t.Status.ExecCmd.Stdout: %v\n", t.Status.ExecCmd.Stdout)
	fmt.Printf("t.Status.ExecCmd.Stderr: %v\n", t.Status.ExecCmd.Stderr)
}

// Execute starts the defined command with it;s arguemnts.
func ExecuteWithTask(taskPointerAddress **Task) {
	l.Debug("ddd")
	t := *taskPointerAddress
	// fmt.Println("In ExecuteWithTask")
	// fmt.Printf("address in main: %p\n", t)
	// fmt.Printf("%+v", t)
	// fmt.Printf("Executing: %s\n", t.Cmd)
	var cancel context.CancelFunc
	if t.Timeout == 0 {
		t.Timeout = DefaultTimeout
	}
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, t.Cmd, t.Args...)
	t.SetEnvironment()

	t.Status.ExecCmd.Stdin = os.Stdin
	t.Status.ExecCmd.Stdout = os.Stdout
	t.Status.ExecCmd.Stderr = os.Stderr
	if err := t.Status.ExecCmd.Run(); err != nil {
		log.Fatal(err)
	}

	l.Debug(fmt.Sprintf("t.Status.ExecCmd.ProcessState.ExitCode(): %v\n", t.Status.ExecCmd.ProcessState.ExitCode()))
	l.Debug(fmt.Sprintf("t.Status.ExecCmd.ProcessState.Pid(): %v\n", t.Status.ExecCmd.ProcessState.Pid()))
	// fmt.Printf("t.Status.ExecCmd.ProcessState.SysUsage(): %v\n", t.Status.ExecCmd.ProcessState.SysUsage())
	// fmt.Printf("t.Status.ExecCmd.Stdout: %v\n", t.Status.ExecCmd.Stdout)
	// fmt.Printf("t.Status.ExecCmd.Stderr: %v\n", t.Status.ExecCmd.Stderr)
}

// SetEnvironment prepares the environment values using pre-defined rules based on regex in the configuration.
func (t *Task) SetEnvironment() {
	var environmentValuesAfterAcceptFilter []string
	var environmentValuesToBeFedToProcess []string
	var environmentValuesConfig []string

	// Traverse the environment values we have and apply the accept filters.

	for _, envKeyValue := range os.Environ() {
		if t.Environment.AcceptFilterRegex != nil {
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
				if CheckRegex(envKeyValue, regex) {
					environmentValuesToBeFedToProcess = append(environmentValuesToBeFedToProcess, envKeyValue)
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

	t.Status.ExecCmd.Env = environmentValuesToBeFedToProcess
}
