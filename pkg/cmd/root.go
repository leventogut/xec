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
	AppName = "butler"
	// ConfigFileNameWithoutExtension is the name of the config file
	ConfigFileNameWithoutExtension = ".butler"
)

var (
	C          xec.Config
	configFile string
	// Verbose defines the verbosity as a boolean.
	Verbose bool = true
	// Debug defines if debug should be enabled.
	Debug bool = false
	// Timeout defines the maximum time the chore execution can take place.
	Timeout int = 10
	// NoColor defines a boolean, when true output will not be colorized.
	NoColor bool = false
	o            = output.NewOutput()

	logFile string = AppName + ".log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "butler [OPTIONS] <alias>",
	Short: "Task executor.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var err error
	o.Output("Available tasks:", "info")
	for _, v := range C.Tasks {
		o.Output(fmt.Sprintf("%s\t%s", v.Alias, v.Description), "info")
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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.butler.yaml)")
	rootCmd.PersistentFlags().BoolVar(&NoColor, "nocolor", true, "Color output. (Default is true i.e. color enabled)")
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", true, "Verbose level output.  (Default is true i.e. verbose output enabled)")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug level output.  (Default is true i.e. debug output enabled)")
	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "", "Filename to use for logging.")

	// Local flags (bare root command):
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	var err error
	viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	viper.SetConfigType("yaml")                         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                            // optionally look for config in the working directory
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

	// Set config
	logFile = AppName + ".log"
	o.SetLogFile(AppName + ".log")
	fmt.Printf("logFile: %s\n", logFile)

	if Debug {
		CJSON, err := json.MarshalIndent(C, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Config in indented JSON:\n %s\n", string(CJSON))
	}

	if Verbose {
		o.Output("Verbose enabled", "debug")
	}
	if Debug {
		//fmt.Printf("Debug enabled\n")
		o.Output("Debug enabled", "debug")
	}
}
