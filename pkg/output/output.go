package output

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type Output interface {
	NoColorFlag() bool
	SetNoColorFlag(b bool)
	LogFile() string
	SetLogFile(f string)
	QuietFlag() bool
	SetQuietFlag(b bool)
	DebugFlag() bool
	SetDebugFlag(b bool)
	VerboseFlag() bool
	SetVerboseFlag(b bool)
	Output()
}
type output struct {
	noColor bool
	logFile string
	writers []io.Writer
	quiet   bool
	debug   bool
	color   *color.Color
	verbose bool
}

var (
	L       *output
	prefix  string
	once    sync.Once
	logFile *os.File
)

func NewOutput(noColor bool, logFile string, quiet bool, debug bool, dev bool, verbose bool) *output {
	return &output{
		noColor: noColor,
		logFile: logFile,
		quiet:   quiet,
		debug:   debug,
		verbose: verbose,
	}
}

func GetInstance() *output {
	once.Do(func() {
		L = NewOutput(false, "", false, false, false, false)
		if !L.quiet {
			L.writers = append(L.writers, os.Stdout)
		}
		// L.setLogFile()
	})
	return L
}
func (o *output) setLogFile() {
	if o.logFile != "" {
		logFile, err := os.OpenFile(L.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Can't open log file:  %v", err)
		}
		o.writers = append(o.writers, logFile)
	}
}
func (o *output) NoColorFlag() bool {
	return o.noColor
}
func (o *output) SetNoColorFlag(b bool) {
	o.noColor = b
}
func (o *output) LogFileFlag() string {
	return o.logFile
}
func (o *output) SetLogFileFlag(f string) {
	o.logFile = f
	o.setLogFile()
}
func (o *output) QuietFlag() bool {
	return o.quiet
}
func (o *output) SetQuietFlag(b bool) {
	o.quiet = b
}
func (o *output) DebugFlag() bool {
	return o.debug
}
func (o *output) SetDebugFlag(b bool) {
	o.debug = b
}
func (o *output) VerboseFlag() bool {
	return o.verbose
}
func (o *output) SetVerboseFlag(b bool) {
	o.verbose = b
}
func (o *output) Debug(m string, values ...interface{}) {
	o.Output(m, "debug", values...)
}

func (o *output) Fatal(m string, values ...interface{}) {
	o.Output(m, "fatal", values...)
	os.Exit(1)
}

func (o *output) Error(m string, values ...interface{}) {
	o.Output(m, "error", values...)
}

func (o *output) Warning(m string, values ...interface{}) {
	o.Output(m, "warning", values...)
}

func (o *output) Info(m string, values ...interface{}) {
	o.Output(m, "info", values...)
}

func (o *output) Success(m string, values ...interface{}) {
	o.Output(m, "success", values...)
}

// Output receives two strings (severity and message and outputs to stdout or
func (o *output) Output(message string, outputType string, values ...interface{}) {
	debugColor := color.FgYellow
	fatalColor := color.FgRed
	errorColor := color.FgRed
	warningColor := color.FgYellow
	infoColor := color.FgHiBlue
	successColor := color.FgGreen

	color.Unset()
	defer color.Unset()

	if outputType == "debug" {
		o.color = color.New(debugColor)
	} else if outputType == "fatal" {
		o.color = color.New(fatalColor)
	} else if outputType == "warning" {
		o.color = color.New(warningColor)
	} else if outputType == "error" {
		o.color = color.New(errorColor)
	} else if outputType == "info" {
		o.color = color.New(infoColor)
	} else if outputType == "success" {
		o.color = color.New(successColor)
	} else {
		o.color.DisableColor()
	}
	if o.noColor {
		o.color.DisableColor()
	}

	w := false
	if !o.quiet {
		if outputType == "error" || outputType == "fatal" {
			w = true
		} else if o.debug {
			w = true
		} else if outputType == "info" && o.verbose {
			w = true
		} else if outputType == "warning" && o.verbose {
			w = true
		} else if outputType == "success" {
			w = true
		} else {
			w = false
		}
	}

	if w {
		now := time.Now().Format(time.RFC3339)

		messageResult := now + " | " + strings.ToUpper(outputType) + " | " + message
		_, err := o.write(fmt.Sprintf(messageResult, values...))
		if err != nil {
			fmt.Printf("Couldn't write log line, error: %+v\n", err)
		}
		if outputType == "fatal" {
			os.Exit(1)
		}
	}
}

func (o *output) write(s string) (n int, err error) {
	if o.logFile != "" {
		defer logFile.Close()
	}
	for _, writer := range o.writers {
		if o.logFile != "" {
			o.color.DisableColor()
		}
		_, err = o.color.Fprintf(writer, s+"\n")
		if err != nil {
			fmt.Printf("err: %v", err)
		}
	}
	return len(s), err
}
