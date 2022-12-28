package output

import (
	"fmt"
	"io"
	"log"
	"os"
)
import "github.com/fatih/color"

type Output interface {
	NoColor() bool
	SetNoColor(nc bool)
	LogFile() string
	SetLogFile(f string)
	Quiet() bool
	SetQuiet(q bool)
}
type output struct {
	noColor bool
	logFile string
	writers []io.Writer
	quiet   bool
	color   *color.Color
}

func NewOutput() *output {
	return &output{
		noColor: false,
		logFile: "",
		quiet:   false,
	}
}

func (o output) NoColor() bool {
	return o.noColor
}
func (o output) SetNoColor(nc bool) {
	o.noColor = nc
}
func (o output) LogFile() string {
	return o.logFile
}
func (o output) SetLogFile(f string) {
	fmt.Println("Setting logFile")
	o.logFile = f
}
func (o output) Quiet() bool {
	return o.quiet
}
func (o output) SetQuiet(q bool) {
	o.quiet = q
}

// Usage:
// o.Output("error", "this is my message", os.Stderr)
// "this is my output candidate".output("warning")

// Output receives two strings (severity and message and outputs to stdout or
func (o output) Output(message string, outputType string, writers ...[]io.Writer) {
	fmt.Printf("%+v", o)
	debugColor := color.FgYellow
	fatalColor := color.FgRed
	errorColor := color.FgRed
	warningColor := color.FgHiYellow
	infoColor := color.FgHiBlue
	successColor := color.FgHiGreen
	defer color.Unset()
	color.Unset()
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
	if !o.quiet {
		if outputType == "error" || outputType == "fatal" {
			o.writers = append(o.writers, os.Stderr)
		} else {
			o.writers = append(o.writers, os.Stdout)
		}
	}
	fmt.Printf("o.logFile: %s\n", o.logFile)
	// Log to file if logFile is set
	if o.logFile != "" {
		f, err := os.OpenFile(o.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Can't open log file:  %v", err)
		}
		defer f.Close()
		o.writers = append(o.writers, f)
	}
	_, err := o.write(message)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
}

func (o output) write(s string) (n int, err error) {
	for _, writer := range o.writers {
		fmt.Printf("writer: %v", writer)
		_, err = o.color.Fprintf(writer, s+"\n")
		if err != nil {
			fmt.Printf("err: %v", err)
		}
	}
	return len(s), err
}
