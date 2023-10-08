package cmd

import (
	"fmt"
	"leventogut/xec/pkg/output"
	"leventogut/xec/pkg/xec"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName                        = "xec"  // AppName is the name of the application
	ConfigFileNameWithoutExtension = ".xec" // ConfigFileNameWithoutExtension is the name of the config file
)

var (
	o               = output.GetInstance()
	C               xec.Config // C is config object.
	ConfigFile      string
	Verbose         bool         // Verbose defines the verbosity as a boolean.
	Debug           bool         // Debug defines if debug should be enabled.
	Dev             bool         // Dev enables development level output
	Timeout         int    = 600 // Timeout defines the maximum time the task execution can take place.
	NoColor         bool         // NoColor defines a boolean, when true output will not be colorized.
	LogFile         string       // Log file name
	Quiet           bool         // Quiet option
	IgnoreErrorFlag bool         // Continue even if the task errors
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xec <flags> <alias> -- [additional-arguments]",
	Short: "Simple command (task) executor.",
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
			if C.TaskDefaults.Timeout == 0 {
				t.Timeout = xec.DefaultTimeout
			}
		}

		o.Success(fmt.Sprintf("IgnoreErrorFlag: +%v\n", IgnoreErrorFlag))
		o.Dev(fmt.Sprintf("t.IgnoreError is %v\n", t.IgnoreError))
		o.Dev(fmt.Sprintf("Execute: C.TaskDefaults.IgnoreError is %v\n", C.TaskDefaults.IgnoreError))

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

		if t.LogFile == "" {
			t.LogFile = C.TaskDefaults.LogFile
		}

		// Add task aliases (sub-commands)
		rootCmd.AddCommand(&cobra.Command{
			Use:   t.Alias,
			Short: t.Description,
			Args:  cobra.ArbitraryArgs,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				o.SetLogFileFlag(LogFile)
				o.SetNoColorFlag(NoColor)
				o.SetQuietFlag(Quiet)
				o.SetDebugFlag(Debug)
				o.SetDevFlag(Dev)
				o.SetVerboseFlag(Verbose)
			},
			Run: func(cmd *cobra.Command, args []string) {
				xec.Execute(&t, args)
			},
		})
	}
	// Traverse TaskLists
	for _, taskListInstance := range C.TaskLists {
		tL := taskListInstance
		// Find tasks from TaskList that matches by alias
		var taskListTasks []*xec.Task
		for _, taskName := range tL.TaskNames {
			// Find tasks pointer address
			for _, tInstance := range C.Tasks {
				t := tInstance
				if taskName == t.Alias {
					if tL.IgnoreError {
						tInstance.IgnoreError = true
					}
					taskListTasks = append(taskListTasks, t)
				}
			}
		}
		// For each TaskList add a command
		rootCmd.AddCommand(&cobra.Command{
			Use:   tL.Alias,
			Short: tL.Description,
			Args:  cobra.ArbitraryArgs,
			Run: func(cmd *cobra.Command, args []string) {
				for _, taskListTask := range taskListTasks {
					// o.Error(fmt.Sprintf("taskListTask: %+v", taskListTask))
					xec.Execute(&taskListTask, taskListTask.Args)
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
	rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is ~/.xec.yaml and/or $PWD/.xec.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Dev, "dev", "z", false, "Enable development level messages.")
	rootCmd.PersistentFlags().BoolVarP(&NoColor, "no-color", "n", false, "Disable color output. (Default is true i.e. color enabled.)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose level output.  (Default is true i.e. verbose output enabled.)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Debug level output.  (Default is true i.e. debug output enabled.)")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "No output.  (Default is false i.e. not quiet.)")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "log-file", "l", "", "Filename to use for logging.")
	rootCmd.PersistentFlags().BoolVarP(&IgnoreErrorFlag, "ignore-error", "i", false, "Ignore error on tasks.")

	viper.BindPFlag("dev", rootCmd.Flags().Lookup("dev"))
	viper.BindPFlag("noColor", rootCmd.Flags().Lookup("no-color"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
	viper.BindPFlag("quiet", rootCmd.Flags().Lookup("quiet"))
	viper.BindPFlag("logFile", rootCmd.Flags().Lookup("log-file"))
	viper.BindPFlag("ignoreErrorFlag", rootCmd.Flags().Lookup("ignore-error"))

	initConfig()
}

func initConfig() {
	var err error
	viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	viper.SetConfigType("yaml")                         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                            // look for config in the working directory
	viper.AddConfigPath("$HOME")                        // look for configs in the $HOME

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}
}
