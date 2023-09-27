package cmd

import (
	"encoding/json"
	"fmt"
	o "leventogut/xec/pkg/output"
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
	C          xec.Config // C is config object.
	ConfigFile string
	Verbose    bool   = true             // Verbose defines the verbosity as a boolean.
	Debug      bool   = false            // Debug defines if debug should be enabled.
	Dev        bool   = false            // Dev enables development level output
	Timeout    int    = 600              // Timeout defines the maximum time the task execution can take place.
	NoColor    bool   = false            // NoColor defines a boolean, when true output will not be colorized.
	LogFile    string = AppName + ".log" // Logfile name
	Quiet             = false
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xec [OPTIONS] <alias>",
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
			t.Timeout = xec.DefaultTimeout
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
		if t.Timeout == 0 {
			if C.TaskDefaults.Timeout != 0 {
				t.Timeout = C.TaskDefaults.Timeout
			} else {
				t.Timeout = xec.DefaultTimeout
			}
		}

		t.Environment.Values = append(t.Environment.Values, C.TaskDefaults.Environment.Values...)

		if t.LogFile == "" {
			t.LogFile = C.TaskDefaults.LogFile
		}

		rootCmd.AddCommand(&cobra.Command{
			Use:   t.Alias,
			Short: t.Description,
			// DisableFlagParsing: true,
			Args: cobra.ArbitraryArgs,
			Run: func(cmd *cobra.Command, args []string) {
				o.L.Debug(fmt.Sprintf("args under sub-commands: %+v\n", args))
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
					taskListTasks = append(taskListTasks, t)
				}
			}
		}
		// For each TaskList add a command
		rootCmd.AddCommand(&cobra.Command{
			Use:   tL.Alias,
			Short: tL.Description,
			// DisableFlagParsing: true,
			Args: cobra.ArbitraryArgs,
			Run: func(cmd *cobra.Command, args []string) {
				for _, taskListTask := range taskListTasks {
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
	rootCmd.PersistentFlags().BoolVar(&Dev, "dev", false, "Enable development level messages.")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "no-color", false, "Disable color output. (Default is true i.e. color enabled.)")
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", true, "Verbose level output.  (Default is true i.e. verbose output enabled.)")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug level output.  (Default is true i.e. debug output enabled.)")
	rootCmd.PersistentFlags().BoolVar(&Quiet, "quiet", false, "No output.  (Default is false i.e. not quiet.)")
	rootCmd.PersistentFlags().StringVar(&LogFile, "log-file", "", "Filename to use for logging.")

	if Dev {
		o.L.Dev(fmt.Sprintf("Dev: %+v", Dev))
		o.L.Dev("Dev enabled")
	}
	if NoColor {
		o.L.Dev(fmt.Sprintf("NoColor: %+v", NoColor))
		o.L.Dev("NoColor enabled")
	}

	if Verbose {
		o.L.Dev(fmt.Sprintf("Verbose: %+v", Verbose))
		o.L.Dev("Verbose enabled")
	}
	if Debug {
		o.L.Dev(fmt.Sprintf("Debug: %+v", Debug))
		o.L.Dev("Debug enabled")
	}
	if Quiet {
		o.L.Dev(fmt.Sprintf("Quiet: %+v", Quiet))
		o.L.Dev("Quiet enabled")
	}

	log.Printf("LogFile->: %+v\n", LogFile)
	o.L.SetLogFileFlag(LogFile)
	o.L.SetNoColorFlag(NoColor)
	o.L.SetQuietFlag(Quiet)
	o.L.SetDebugFlag(Debug)
	log.Printf("Dev: %+v\n", Dev)
	o.L.SetDevFlag(Dev)
	o.L.SetLogFileFlag(LogFile)

	initConfig()
}

func initConfig() {
	var err error
	viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	viper.SetConfigType("yaml")                         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                            // look for config in the working directory
	viper.AddConfigPath("$HOME")                        // look for configs in the $HOME
	// defaults
	viper.SetDefault("Verbose", Verbose)
	viper.SetDefault("Debug", Debug)
	viper.SetDefault("Timeout", Timeout)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}

	// viper.Debug()
	if Dev {
		CJSON, err := json.MarshalIndent(C, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		o.L.Debug(fmt.Sprintf("Config in indented JSON:\n %s\n", string(CJSON)))
	}
}
