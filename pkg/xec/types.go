package xec

import (
	"context"
	"os/exec"
)

// Config represents the configuration read from file or env values.
type Config struct {
	Verbose      bool         `yaml:"verbose" json:"verbose"`
	Debug        bool         `yaml:"debug" json:"debug"`
	NoColor      bool         `yaml:"no-color" json:"no-color"`
	Quiet        bool         `yaml:"quiet" json:"quiet"`
	LogFile      string       `yaml:"logFile" json:"logFile"`
	LogDir       string       `yaml:"logDir" json:"logDir"`
	RestartLimit int          `yaml:"restartLimit" json:"restartLimit"`
	Imports      []string     `yaml:"imports" json:"imports"`
	TaskDefaults TaskDefaults `yaml:"taskDefaults" json:"taskDefaults"`
	Tasks        []*Task      `yaml:"tasks" json:"tasks"`
	TaskLists    []*TaskList  `yaml:"taskLists" json:"taskLists"`
	Path         string
	Namespace    string `yaml:"namespace" json:"namespace"`
}

// TaskDefaults defines the default values for all tasks mentioned.
type TaskDefaults struct {
	Timeout      int         `yaml:"timeout" json:"timeout"`
	RestartLimit int         `yaml:"restartLimit" json:"restartLimit"`
	Environment  Environment `yaml:"environment" json:"environment"`
	LogFile      string      `yaml:"logFile" json:"logFile"`
	IgnoreError  bool        `yaml:"ignoreError" json:"ignoreError"`
}

// Task is a combination of alias, description, command, and arguments.
type Task struct {
	Alias            string      `yaml:"alias" json:"alias"`
	Description      string      `yaml:"description" json:"description"`
	Cmd              string      `yaml:"cmd" json:"cmd"`
	Args             []string    `yaml:"args" json:"args"`
	ExtraArgs        []string    // Extra args coming from the command line not the configuration.
	Timeout          int         `yaml:"timeout" json:"timeout"`
	Environment      Environment `yaml:"environment" json:"environment"`
	LogFile          string      `yaml:"logFile" json:"logFile"`
	LogDir           string      `yaml:"logDir" json:"logDir"`
	IgnoreError      bool        `yaml:"ignoreError" json:"ignoreError"`
	RestartOnSuccess bool        `yaml:"restartOnSuccess" json:"restartOnSuccess"`
	RestartOnFailure bool        `yaml:"restartOnFailure" json:"restartOnFailure"`
	RestartLimit     int         `yaml:"restartLimit" json:"restartLimit"`
	NumberOfRestarts int
	Directory        string     `yaml:"directory" json:"directory"`
	Status           TaskStatus `yaml:"taskStatus" json:"taskStatus"`
}

// Environment defines the environment key/values that shoul be feed to the process.
type Environment struct {
	PassOn            bool                `yaml:"passOn" json:"passOn"`
	Values            []map[string]string `yaml:"values" json:"values"`
	AcceptFilterRegex []string            `yaml:"acceptFilterRegex" json:"acceptFilterRegex"`
	RejectFilterRegex []string            `yaml:"rejectFilterRegex" json:"rejectFilterRegex"`
}

// TaskStatus represents a chore's status and properties of underlying exec structure.
type TaskStatus struct {
	Started     bool
	Finished    bool
	Success     bool
	ExitCode    int
	ExecContext context.Context
	ExecCmd     *exec.Cmd
}

// TaskList is a virtual list of tasks to be run in the given order.
type TaskList struct {
	Alias       string   `yaml:"alias" json:"alias"`
	Description string   `yaml:"description" json:"description"`
	TaskAliases []string `yaml:"taskAliases" json:"taskAliases"`
	LogFile     string   `yaml:"logFile" json:"logFile"`
	IgnoreError bool     `yaml:"ignoreError" json:"ignoreError"`
	Parallel    bool
}
