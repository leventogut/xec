package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"leventogut/xec/pkg/xec"

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
	// Verbose defines the verbosity as a boolean.
	Verbose bool
	// Debug defines if debug should be enabled.
	Debug bool
	// Timeout defines the maximum time the chore execution can take place.
	Timeout int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "butler [OPTIONS] <alias>",
	Short: "Task finisher.",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if Debug {

		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigName(ConfigFileNameWithoutExtension) // name of config file (without extension)
	viper.SetConfigType("yaml")                         // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")                            // optionally look for config in the working directory
	err := viper.ReadInConfig()                         // Find and read the config file
	if err != nil {                                     // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	// defaults
	viper.SetDefault("Verbose", false)
	viper.SetDefault("Timeout", 30)

	// environment values

	// Set config
	Verbose = viper.GetBool("Verbose")
	Timeout = viper.GetInt("Timeout")

	// Global flags:
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/.xec.yaml)")

	// Local flags (bare root command):
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	var C xec.Config
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatalf("Can't decode config, error: %v", err)
	}
	if Debug {
		CJSON, err := json.MarshalIndent(C, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Config in indented JSON:\n %s\n", string(CJSON))
	}

	for _, v := range C.Tasks {
		fmt.Printf("%v", v)
	}
}
