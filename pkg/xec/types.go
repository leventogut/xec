package xec

import (
	"context"
	"os/exec"
)

// Config represents the configuration read from file or env values.
type Config struct {
	Verbose     bool        `yaml:"verbose" json:"verbose"`
	Debug       bool        `yaml:"debug" json:"debug"`
	Environment Environment `yaml:"environment" json:"environment"`
	Tasks       []Task      `yaml:"tasks" json:"tasks"`
}

// Environment defines the environment key/values that shoul be feed to the process.
type Environment struct {
	PassOn            bool              `yaml:"passOn" json:"passOn"`
	Values            map[string]string `yaml:"values" json:"values"`
	AcceptFilterRegex []string          `yaml:"acceptFilterRegex" json:"acceptFilterRegex"`
	RejectFilterRegex []string          `yaml:"rejectFilterRegex" json:"rejectFilterRegex"`
}

// Task is a combination of alias, description, command, and arguments.
type Task struct {
	Alias       string              `yaml:"alias" json:"alias"`
	Description string              `yaml:"description" json:"description"`
	Cmd         string              `yaml:"cmd" json:"cmd"`
	Args        []string            `yaml:"args" json:"args"`
	Labels      []string            `yaml:"labels" json:"labels"`
	Tags        []map[string]string `yaml:"tags" json:"tags"`
	Timeout     int                 `yaml:"timeout" json:"timeout"`
	Environment Environment         `yaml:"environment" json:"environment"`
	LogFile     string              `yaml:"logFile" json:"logFile"`
	Status      TaskStatus
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
