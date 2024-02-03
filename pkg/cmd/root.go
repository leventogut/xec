package cmd

import (
	"flag"
	"fmt"
	"leventogut/xec/pkg/output"
	"leventogut/xec/pkg/xec"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName                        = "xec"  // AppName is the name of the application
	ConfigFileNameWithoutExtension = ".xec" // ConfigFileNameWithoutExtension is the name of the config file
)

var (
	o               = output.GetInstance()
	C               xec.Config       // C is config object.
	ConfigFile      string           // Custom configuration file.
	Verbose         bool             // Verbose defines the verbosity as a boolean.
	Debug           bool             // Debug defines if debug should be enabled.
	Dev             bool             // Dev enables development level output
	Timeout         int        = 600 // Timeout defines the maximum time the task execution can take place.
	NoColor         bool             // NoColor defines a boolean, when true output will not be colorized.
	LogFile         string           // Log file name
	Quiet           bool             // Quiet option
	IgnoreErrorFlag bool             // Continue even if the task errors
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xec <flags> <alias> -- [args-to-be-passed]",
	Short: "Simple command executor.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

// Execute adds all child commands (task aliases) to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	for _, tInstance := range C.Tasks {
		t := tInstance

		if t.Timeout == 0 {
			t.Timeout = C.TaskDefaults.Timeout
			if t.Timeout == 0 {
				t.Timeout = xec.DefaultTimeout
			}
		}

		if IgnoreErrorFlag {
			t.IgnoreError = true
		}
		if !t.IgnoreError {
			t.IgnoreError = C.TaskDefaults.IgnoreError
		}
		// Copy defaults from main config/taskDefaults to the task if task value is empty / undefined.
		if !t.Environment.PassOn {
			t.Environment.PassOn = C.TaskDefaults.Environment.PassOn
		}

		if t.Environment.AcceptFilterRegex == nil {
			t.Environment.AcceptFilterRegex = C.TaskDefaults.Environment.AcceptFilterRegex
		}
		if t.Environment.RejectFilterRegex == nil {
			t.Environment.RejectFilterRegex = C.TaskDefaults.Environment.RejectFilterRegex
		}

		t.Environment.Values = append(t.Environment.Values, C.TaskDefaults.Environment.Values...)

		if t.Timeout == 0 {
			if C.TaskDefaults.Timeout != 0 {
				t.Timeout = C.TaskDefaults.Timeout
			} else {
				t.Timeout = xec.DefaultTimeout
			}
		}

		// --log-file overrides taskDefaults and task logFile, task overrides taskDefaults
		if t.LogFile == "" {
			t.LogFile = C.TaskDefaults.LogFile
		}

		if t.LogFile == "auto" {
			now := time.Now().Format(time.RFC3339Nano)
			t.LogFile = "xec-log-" + t.Alias + "-" + now + ".log"
			o.SetLogFileFlag(t.LogFile)
		}

		// Add task aliases (sub-commands)
		rootCmd.AddCommand(&cobra.Command{
			Use:   t.Alias,
			Short: t.Cmd + strings.Join(t.Args[:], " "),
			Long:  t.Description,
			Args:  cobra.ArbitraryArgs,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				o.SetLogFileFlag(t.LogFile)
				o.SetNoColorFlag(NoColor)
				o.SetQuietFlag(Quiet)
				o.SetDebugFlag(Debug)
				o.SetDevFlag(Dev)
				o.SetVerboseFlag(Verbose)
			},
			Run: func(cmd *cobra.Command, args []string) {
				xec.Execute(&t)
			},
		})
	}
	// Traverse TaskLists
	for _, taskListInstance := range C.TaskLists {
		tL := taskListInstance
		// Find tasks from TaskList that matches by alias
		var taskListTasks []*xec.Task
		for _, taskName := range tL.TaskAliases {
			// Find tasks pointer address
			for _, tInstance := range C.Tasks {
				t := tInstance
				if taskName == t.Alias {
					if tL.IgnoreError {
						tInstance.IgnoreError = true
					}
					if tL.LogFile != "" {
						tInstance.LogFile = tL.LogFile
					}
					if tInstance.LogFile == "" {
						tInstance.LogFile = C.TaskDefaults.LogFile
					}
					taskListTasks = append(taskListTasks, t)
				}
			}
		}

		if tL.LogFile == "auto" {
			now := time.Now().Format(time.RFC3339Nano)
			tL.LogFile = "xec-log-" + tL.Alias + "-" + now + ".log"
			o.SetLogFileFlag(tL.LogFile)
		}

		// For each TaskList add a command
		rootCmd.AddCommand(&cobra.Command{
			Use:   tL.Alias,
			Short: tL.Description,
			Args:  cobra.ArbitraryArgs,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				o.SetLogFileFlag(tL.LogFile)
				o.SetNoColorFlag(NoColor)
				o.SetQuietFlag(Quiet)
				o.SetDebugFlag(Debug)
				o.SetDevFlag(Dev)
				o.SetVerboseFlag(Verbose)
			},
			Run: func(cmd *cobra.Command, args []string) {
				if tL.Parallel {
					var wg sync.WaitGroup
					for _, taskListTask := range taskListTasks {
						taskListTask := taskListTask
						wg.Add(1)
						go xec.ExecuteWithWaitGroups(&wg, &taskListTask)
					}
					wg.Wait()

				} else {
					for _, taskListTask := range taskListTasks {
						xec.Execute(&taskListTask)
					}
				}
			},
		})
	}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags:
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "", "", "config file to read (default is ~/.xec.yaml,  $PWD/.xec.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&NoColor, "no-color", "", false, "Disable color output.")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "", false, "Verbose level output.")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "Debug level output.")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "", false, "No output except errors].")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "log-file", "", "", "Filename to use for logging.")
	rootCmd.PersistentFlags().BoolVarP(&IgnoreErrorFlag, "ignore-error", "", false, "Ignore errors on tasks.")

	flag.StringVar(&ConfigFile, "config", "", "")
	flag.BoolVar(&NoColor, "no-color", false, "")
	flag.BoolVar(&Verbose, "verbose", false, "")
	flag.BoolVar(&Debug, "debug", false, "")
	flag.BoolVar(&Quiet, "quiet", false, "")
	flag.StringVar(&LogFile, "log-file", "", "")
	flag.BoolVar(&IgnoreErrorFlag, "ignore-error", false, "")

	viper.BindPFlags(rootCmd.PersistentFlags())

	initConfig()
}

func initConfig() {
	var err error

	flag.Parse()

	if ConfigFile != "" {
		viper.SetConfigName(ConfigFile)

	} else {
		viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	}

	viper.SetConfigType("yaml")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")     // look for config in the working directory
	viper.AddConfigPath("$HOME") // look for configs in the $HOME

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}
}
