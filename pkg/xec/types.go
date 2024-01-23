package xec

import (
	"context"
	"os/exec"
)

// Config represents the configuration read from file or env values.
type Config struct {
	TaskDefaults TaskDefaults `yaml:"taskDefaults" json:"taskDefaults"`
	Tasks        []*Task      `yaml:"tasks" json:"tasks"`
	TaskLists    []*TaskList  `yaml:"taskLists" json:"taskLists"`
}

// TaskDefaults defines the default values for all tasks mentioned.
type TaskDefaults struct {
	Debug       bool        `yaml:"debug" json:"debug"`
	Timeout     int         `yaml:"timeout" json:"timeout"`
	Environment Environment `yaml:"environment" json:"environment"`
	LogFile     string      `yaml:"logFile" json:"logFile"`
	IgnoreError bool        `yaml:"ignoreError" json:"ignoreError"`
}

// Task is a combination of alias, description, command, and arguments.
type Task struct {
	Alias       string      `yaml:"alias" json:"alias"`
	Description string      `yaml:"description" json:"description"`
	Cmd         string      `yaml:"cmd" json:"cmd"`
	Args        []string    `yaml:"args" json:"args"`
	ExtraArgs   []string    // Extra args coming from the command line not the configuration.
	Timeout     int         `yaml:"timeout" json:"timeout"`
	Environment Environment `yaml:"environment" json:"environment"`
	LogFile     string      `yaml:"logFile" json:"logFile"`
	IgnoreError bool        `yaml:"ignoreError" json:"ignoreError"`
	Status      TaskStatus  `yaml:"taskStatus" json:"taskStatus"`
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
	IgnoreError bool     `yaml:"ignoreError" json:"ignoreError"`
	Parallel    bool
}
