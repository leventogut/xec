package cmd

import (
	"encoding/json"
	"fmt"
	"leventogut/xec/pkg/output"
	"leventogut/xec/pkg/xec"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// AppName is the name of the application
	AppName = "xec"
	// ConfigFileNameWithoutExtension is the name of the config file
	ConfigFileNameWithoutExtension = "xec"
)

var (
	C          xec.Config
	configFile string
	// Verbose defines the verbosity as a boolean.
	Verbose bool = true
	// Debug defines if debug should be enabled.
	Debug bool = true
	// Timeout defines the maximum time the chore execution can take place.
	Timeout int = 10
	// NoColor defines a boolean, when true output will not be colorized.
	NoColor bool   = false
	LogFile string = AppName + ".log"
	Quiet          = false
	// TODO import from outpiut package and set the values accordingly with set functions.
	l = output.Log
	// Log = output.NewOutput(NoColor, LogFile, Quiet, Debug)
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xec [OPTIONS] <alias>",
	Short: "Task executor.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	for _, t := range C.Tasks {
		l.Info(fmt.Sprintf("%s\t%s\t%s\t%s", t.Alias, t.Description, t.Cmd, t.Args))
		rootCmd.AddCommand(&cobra.Command{
			Use:   t.Alias,
			Short: t.Description,
			Long:  t.Description,
			Run: func(cmd *cobra.Command, args []string) {
				xec.ExecuteWithTask(&t)
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
	cobra.OnInitialize(initConfig)
	// Global flags:
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $PWD/.xec.yaml)")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "nocolor", true, "Color output. (Default is true i.e. color enabled)")
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", true, "Verbose level output.  (Default is true i.e. verbose output enabled)")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug level output.  (Default is true i.e. debug output enabled)")
	rootCmd.PersistentFlags().StringVar(&LogFile, "Logfile", "", "Filename to use for logging.")

	// Local flags (bare root command):
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initConfig()
}

func initConfig() {
	var err error
	viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	viper.SetConfigType("yaml")                         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                            // look for config in the working directory
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

	if Debug {
		CJSON, err := json.MarshalIndent(C, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		l.Debug(fmt.Sprintf("Config in indented JSON:\n %s\n", string(CJSON)))
	}

	// if Verbose {
	// 	o.Output("Verbose enabled", "debug")
	// }
	if Debug {
		//fmt.Printf("Debug enabled\n")
		l.Debug("Debug enabled")
	}
}
