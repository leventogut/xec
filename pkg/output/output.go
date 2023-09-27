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

type Output interface {
	NoColorFlag() bool
	SetNoColorFlag(b bool)
	LogFile() string
	SetLogFile(f string)
	QuietFlag() bool
	SetQuietFlag(b bool)
	DevFlag() bool
	SetDevFlag(b bool)
	DebugFlag() bool
	SetDebugFlag(b bool)
	VerboseFlag() bool
	SetVerboseFlag(b bool)
}
type output struct {
	noColor bool
	logFile string
	writers []io.Writer
	quiet   bool
	debug   bool
	color   *color.Color
	dev     bool
	verbose bool
}

var (
	L      *output
	prefix string
)

func init() {
	L = NewOutput(false, "", false, false, false, false)
	// L = new(output)
	L.SetNoColorFlag(false)
	L.SetDebugFlag(false)
	L.SetDevFlag(false)
	L.SetLogFileFlag("")
	L.SetQuietFlag(false)
	L.SetVerboseFlag(false)
}

func NewOutput(noColor bool, logFile string, quiet bool, debug bool, dev bool, verbose bool) *output {
	log.Println("output initialized")
	return &output{
		noColor: noColor,
		logFile: logFile,
		quiet:   quiet,
		debug:   debug,
		dev:     dev,
		verbose: verbose,
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
}
func (o *output) QuietFlag() bool {
	return o.quiet
}
func (o *output) SetQuietFlag(b bool) {
	o.quiet = b
}
func (o *output) DevFlag() bool {
	return o.dev
}
func (o *output) SetDevFlag(b bool) {
	o.dev = b
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
func (o *output) Dev(m string) {
	o.Output(m, "dev")
}
func (o *output) Debug(m string) {
	o.Output(m, "debug")
}

func (o *output) Fatal(m string) {
	o.Output(m, "fatal")
	os.Exit(1)
}

func (o *output) Error(m string) {
	o.Output(m, "error")
}

func (o *output) Warning(m string) {
	o.Output(m, "warning")
}

func (o *output) Info(m string) {
	o.Output(m, "info")
}

func (o *output) Success(m string) {
	o.Output(m, "success")
}

// Output receives two strings (severity and message and outputs to stdout or
func (o output) Output(message string, outputType string, writers ...[]io.Writer) {
	devColor := color.FgHiCyan
	debugColor := color.FgYellow
	fatalColor := color.FgRed
	errorColor := color.FgRed
	warningColor := color.FgHiYellow
	infoColor := color.FgHiBlue
	successColor := color.FgHiGreen
	defer color.Unset()
	color.Unset()
	if outputType == "dev" {
		o.color = color.New(devColor)
	} else if outputType == "debug" {
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
	now := time.Now().Format(time.RFC3339)
	if !o.quiet {
		log.Println("not quiet")
		if outputType == "error" || outputType == "fatal" {
			o.writers = append(o.writers, os.Stderr)
		} else if outputType == "debug" && o.debug {
			o.writers = append(o.writers, os.Stdout)
			prefix = now + " | "
		} else if outputType == "info" || outputType == "warning" && o.verbose {
			o.writers = append(o.writers, os.Stdout)
		} else if outputType == "dev" && o.dev {
			o.writers = append(o.writers, os.Stdout)
			prefix = now + " | "
		} else if outputType == "success" {
			o.writers = append(o.writers, os.Stdout)
		} else {
			o.writers = append(o.writers, os.Stdout)
			prefix = "UNDEFINED | "
		}
	}
	// Log to file if logFile is set
	if o.logFile != "" {
		f, err := os.OpenFile(o.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Can't open log file:  %v", err)
		}
		defer f.Close()
		o.writers = append(o.writers, f)
	}
	_, err := o.write("[" + strings.ToUpper(outputType) + "]" + " | " + message)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}

func (o output) write(s string) (n int, err error) {
	for _, writer := range o.writers {
		_, err = o.color.Fprintf(writer, prefix+s+"\n")
		if err != nil {
			fmt.Printf("err: %v", err)
		}
	}
	return len(s), err
}
