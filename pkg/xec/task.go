package xec

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var ()

// Execute starts the defined command with it;s arguemnts.
func (t Task) Execute() {
	cmdInput := t.Cmd
	fmt.Printf("Executing: %v %v\n", cmdInput, t.Args)
	var cancel context.CancelFunc
	t.Status.ExecContext, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout)*time.Second)
	defer cancel()
	t.Status.ExecCmd = exec.CommandContext(t.Status.ExecContext, cmdInput, t.Args...)
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

// SetEnvironment prepares the environment values using pre-defined rules based on regex in the configuration.
func (t Task) SetEnvironment() {
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
	for _, envKeyValue := range c.Environment.Values {
		environmentValuesConfig = append(environmentValuesConfig, envKeyValue)
	}
	environmentValuesToBeFedToProcess = append(environmentValuesToBeFedToProcess, environmentValuesConfig...)

	t.Status.ExecCmd.Env = environmentValuesToBeFedToProcess
}
