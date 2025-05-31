package output

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Output struct {
	noColor bool
	logFile string
	writers []io.Writer
	quiet   bool
	debug   bool
	color   *color.Color
	verbose bool
}

var (
	L       *Output
	logFile *os.File
)

func NewOutput(noColor bool, logFile string, quiet bool, debug bool, dev bool, verbose bool) *Output {
	return &Output{
		noColor: noColor,
		logFile: logFile,
		quiet:   quiet,
		debug:   debug,
		verbose: verbose,
	}
}

// GetInstance returns a new instance of Output, each task has it's own instance.
func GetInstance() *Output {
	L = NewOutput(false, "", false, false, false, false)
	if !L.quiet {
		L.writers = append(L.writers, os.Stdout)
	}
	return L
}
func (o *Output) setLogFile() {
	if o.logFile != "" {
		logFile, err := os.OpenFile(L.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Can't open log file:  %v", err)
		}
		o.writers = append(o.writers, logFile)
	}
}
func (o *Output) NoColorFlag() bool {
	return o.noColor
}
func (o *Output) SetNoColorFlag(b bool) {
	o.noColor = b
}
func (o *Output) LogFileFlag() string {
	return o.logFile
}
func (o *Output) SetLogFileFlag(f string) {
	o.logFile = f
	o.setLogFile()
}
func (o *Output) QuietFlag() bool {
	return o.quiet
}
func (o *Output) SetQuietFlag(b bool) {
	o.quiet = b
}
func (o *Output) DebugFlag() bool {
	return o.debug
}
func (o *Output) SetDebugFlag(b bool) {
	o.debug = b
}
func (o *Output) VerboseFlag() bool {
	return o.verbose
}
func (o *Output) SetVerboseFlag(b bool) {
	o.verbose = b
}
func (o *Output) Debug(m string, values ...interface{}) {
	o.WriteOutput(m, "debug", values...)
}

func (o *Output) Fatal(m string, values ...interface{}) {
	o.WriteOutput(m, "fatal", values...)
	os.Exit(1)
}

func (o *Output) Error(m string, values ...interface{}) {
	o.WriteOutput(m, "error", values...)
}

func (o *Output) Warning(m string, values ...interface{}) {
	o.WriteOutput(m, "warning", values...)
}

func (o *Output) Normal(m string, values ...interface{}) {
	o.WriteOutput(m, "normal", values...)
}
func (o *Output) Info(m string, values ...interface{}) {
	o.WriteOutput(m, "info", values...)
}

func (o *Output) Success(m string, values ...interface{}) {
	o.WriteOutput(m, "success", values...)
}

// Output receives two strings (severity and message and outputs to stdout or
func (o *Output) WriteOutput(message string, outputType string, values ...interface{}) {
	debugColor := color.FgYellow
	fatalColor := color.FgRed
	errorColor := color.FgRed
	warningColor := color.FgYellow
	normalColor := color.FgWhite
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
	} else if outputType == "normal" {
		o.color = color.New(normalColor)
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
		if outputType == "fatal" {
			w = true
		} else if outputType == "error" && o.verbose {
			w = true
		} else if o.debug {
			w = true
		} else if outputType == "info" && o.verbose {
			w = true
		} else if outputType == "warning" && o.verbose {
			w = true
		} else if outputType == "success" && o.verbose {
			w = true
		} else if outputType == "normal" {
			w = true
		} else {
			w = false
		}
	}

	if w {
		now := time.Now().Format(time.RFC3339)
		var messageResult string
		if outputType == "normal" {
			messageResult = message
		} else {
			messageResult = now + " | " + strings.ToUpper(outputType) + " | " + message
		}

		_, err := o.write(fmt.Sprintf(messageResult, values...))
		if err != nil {
			fmt.Printf("Couldn't write log line, error: %+v\n", err)
		}
		if outputType == "fatal" {
			os.Exit(1)
		}
	}
}

func (o *Output) write(s string) (n int, err error) {
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
