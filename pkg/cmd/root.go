package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/leventogut/xec/pkg/output"
	"github.com/leventogut/xec/pkg/xec"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName                                      = "xec"  // AppName is the name of the application
	DefaultConfigFileNameWithoutExtension        = ".xec" // DefaultConfigFileNameWithoutExtension is the name of the config file
	DefaultConfigExtension                string = "yaml" // DefaultConfigExtension is the extension used by default
)

var (
	o                 = output.GetInstance()
	Configs           []*xec.Config
	ConfigFileFlag    string = ""    // Custom configuration file.
	Verbose           bool   = true  // Verbose defines the verbosity as a boolean.
	Debug             bool   = false // Debug defines if debug should be enabled.
	Quiet             bool   = false // Quiet option nes the maximum time the task execution can take place.
	NoColor           bool   = false // NoColor defines a boolean, when true output will not be colorized.
	LogDir            string = ""    // LogDir is the destination directory for log files.
	LogFile           string = ""    // Log file name
	IgnoreErrorFlag   bool   = false // Continue even if the task errors
	Timeout           int    = 600   // Timeout for tasks' execution context.
	tLogFile          string = ""
	InitConfiguration string = `# yaml-language-server: $schema=https://raw.githubusercontent.com/leventogut/xec/main/schema/xec-yaml-schema.json
tasks:
  - alias: myCommand
    cmd: echo
    args:
      - "my command is run"
    description: run it via xec myCommand
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xec <flags> <alias> -- [additional-args]",
	Short: "Simple command executor.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(1)
		}
	},
}

// Execute adds all child commands (task aliases) to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	o.SetVerboseFlag(Verbose)
	o.SetDebugFlag(Debug)
	o.SetQuietFlag(Quiet)
	o.SetNoColorFlag(NoColor)

	for _, C := range Configs {
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

			if t.LogFile != "" {

			} else if C.TaskDefaults.LogFile != "" {
				t.LogFile = C.TaskDefaults.LogFile
			} else if C.LogFile != "" {
				t.LogFile = C.LogFile
			}

			if t.LogFile == "auto" {
				now := time.Now().Format(time.RFC3339Nano)
				t.LogFile = "xec-log-" + t.Alias + "-" + now + ".log"
			}

			if t.LogFile != "" {
				t.LogFile = LogDir + t.LogFile
			}

			// Add task aliases (sub-commands)
			rootCmd.AddCommand(&cobra.Command{
				Use:   t.Alias,
				Short: t.Cmd + " " + strings.Join(t.Args[:], " "),
				Long:  t.Description,
				Args:  cobra.ArbitraryArgs,
				PersistentPreRun: func(cmd *cobra.Command, args []string) {
					if t.LogFile != "" {
						o.SetLogFileFlag(t.LogFile)
					}
				},
				Run: func(cmd *cobra.Command, args []string) {
					executeTask := func() {
						xec.Execute(&t)
					}
					_ = spinner.New().
						Title("Task " + t.Alias + " is running.").
						Action(executeTask).
						Run()
				},
			})
		}
		// Traverse TaskLists
		for _, taskListInstance := range C.TaskLists {
			tL := taskListInstance
			// Find tasks from TaskList that matches by alias
			var taskListTasks []*xec.Task
			for _, taskName := range tL.TaskAliases {
				for _, tInstance := range C.Tasks {
					t := tInstance
					if taskName == t.Alias {
						// Properties to be transferred from task list to task
						if tL.IgnoreError {
							t.IgnoreError = true
						}
						if tL.LogFile != "" {
						} else if C.TaskDefaults.LogFile != "" {
							tL.LogFile = C.TaskDefaults.LogFile
						} else if C.LogFile != "" {
							tL.LogFile = C.LogFile
						}

						if tL.LogFile == "auto" {
							now := time.Now().Format(time.RFC3339Nano)
							tL.LogFile = "xec-log-" + tL.Alias + "-" + now + ".log"
						}

						t.LogFile = LogDir + tL.LogFile
						taskListTasks = append(taskListTasks, t)
					}
				}
			}
			tL.LogFile = LogDir + tL.LogFile

			// For each TaskList add a command
			rootCmd.AddCommand(&cobra.Command{
				Use:   tL.Alias,
				Short: tL.Description,
				Args:  cobra.ArbitraryArgs,
				PersistentPreRun: func(cmd *cobra.Command, args []string) {
					if tL.LogFile != "" {
						o.SetLogFileFlag(tL.LogFile)
					}

				},
				Run: func(cmd *cobra.Command, args []string) {
					taskListStartTime := time.Now()
					o.Info("Task list %+v is starting.", tL.Alias)
					o.Info("Task list %+v is logged to %+v.", tL.Alias, tL.LogFile)
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

							executeTask := func() {
								xec.Execute(&taskListTask)
							}
							_ = spinner.New().
								Title("Task " + taskListTask.Alias + " is running.").
								Action(executeTask).
								Run()
						}
					}
					taskListFinishTime := time.Now()
					taskDuration := taskListFinishTime.Sub(taskListStartTime)
					o.Success("Task list " + tL.Alias + " finished in " + taskDuration.String() + ".")
				},
			})
		}
	}
	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags:
	rootCmd.PersistentFlags().StringVarP(&ConfigFileFlag, "config", "", "", "config file to read (default is ~/.xec.yaml,  $PWD/.xec.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&NoColor, "no-color", "", false, "Disable color output.")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "", false, "Verbose level output.")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "Debug level output.")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "", false, "No output except errors].")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "log-file", "", "", "Filename to use for logging.")
	rootCmd.PersistentFlags().StringVarP(&LogDir, "log-dir", "", "", "Directory to use for logging.")
	rootCmd.PersistentFlags().BoolVarP(&IgnoreErrorFlag, "ignore-error", "", false, "Ignore errors on tasks.")

	// Flag package is used due to the fact that Cobra loads flag values quite late.
	// Only reading values and updating global vars is used.
	flag.StringVar(&ConfigFileFlag, "config", "", "")
	flag.BoolVar(&NoColor, "no-color", false, "")
	flag.BoolVar(&Verbose, "verbose", false, "")
	flag.BoolVar(&Debug, "debug", false, "")
	flag.BoolVar(&Quiet, "quiet", false, "")
	flag.StringVar(&LogFile, "log-file", "", "")
	flag.StringVar(&LogDir, "log-dir", "", "")
	flag.BoolVar(&IgnoreErrorFlag, "ignore-error", false, "")

	_ = viper.BindPFlags(rootCmd.PersistentFlags())

	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "initialize a configuration file in the current directory.",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			initConfigFile("./" + DefaultConfigFileNameWithoutExtension + "." + DefaultConfigExtension)
		},
	})

	initConfig()

	if Debug {
		Verbose = true
	}
}

func initConfig() {
	// Parse flags
	flag.Parse()

	if ConfigFileFlag != "" {
		extension := strings.TrimPrefix(filepath.Ext(ConfigFileFlag), ".")
		configName := strings.TrimSuffix(filepath.Base(ConfigFileFlag), filepath.Ext(ConfigFileFlag))
		path := filepath.Dir(ConfigFileFlag)
		CreateViperInstanceFromConfig(path, configName, extension)
	}
	// home directory
	CreateViperInstanceFromConfig("$HOME", DefaultConfigFileNameWithoutExtension, DefaultConfigExtension)

	// current working directory
	CreateViperInstanceFromConfig(".", DefaultConfigFileNameWithoutExtension, DefaultConfigExtension)

	if len(Configs) < 1 {
		o.Warning("No configuration file found.")

		o.Warning("You can generate a skeleton configuration via: %s init", AppName)
	}
}

func CreateViperInstanceFromConfig(path string, configName string, extension string) {
	var err error
	if extension == "" {
		extension = DefaultConfigExtension
	}

	viperInstance := viper.New()
	viperInstance.SetConfigName(configName) // name of config file (without extension)
	viperInstance.SetConfigType(extension)  // REQUIRED if the config file does not have the extension in the name
	viperInstance.AddConfigPath(path)       // look for config in the working directory

	o.Debug("Trying to read config file, [%v/%v.%v].", path, configName, extension)

	err = viperInstance.ReadInConfig() // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		o.Warning("Can't read config file, [%v/%v.%v], skipping", path, configName, extension)
		return
	}

	o.Success("Loaded config file, [%v/%v.%v]", path, configName, extension)

	var config *xec.Config

	err = viperInstance.Unmarshal(&config)
	Configs = append(Configs, config)

	// Root variables form config
	if config.Verbose {
		Verbose = true
	}
	if config.Debug {
		Debug = true
	}
	if config.NoColor {
		NoColor = true
	}
	if config.Quiet {
		Quiet = true
	}
	if config.LogDir != "" {
		LogDir = strings.TrimSuffix(config.LogDir, "/") + "/"
	}
	for _, importConfig := range config.Imports {
		importConfigConfigName := strings.TrimSuffix(filepath.Base(importConfig), filepath.Ext(importConfig))
		importConfigExtension := strings.TrimPrefix(filepath.Ext(importConfig), ".")
		importConfigPath := filepath.Dir(importConfig)
		CreateViperInstanceFromConfig(importConfigPath, importConfigConfigName, importConfigExtension)
	}
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}
}

// initConfigFile checks the given path for configuration file.
// If the file doesn't exist, creates it.
func initConfigFile(fileName string) {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		o.Info("Configuration file [%+v] not found.\n", fileName)
		configFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			o.Fatal("Can't open configuration file: %v, error:  %v\n", fileName, err)
		}
		_, err = configFile.WriteString(InitConfiguration)
		if err != nil {
			o.Fatal("Can't write to file, error:  %v\n", err)
		}
		o.Success("Init configuration is written to file %v.", fileName)
	}
}
