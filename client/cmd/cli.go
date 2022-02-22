package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/geoffjay/plantd/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type server struct {
	Endpoint string `mapstructure:"endpoint"`
}

type clientConfig struct {
	Server server `mapstructure:"server"`
}

var (
	cfgFile  string
	config   clientConfig
	endpoint string
	Verbose  bool

	cliCmd = &cobra.Command{
		Use:   "plant",
		Short: "Application to control plantd services",
		Long:  `A control utility for interacting with plantd services.`,
	}
)

func Execute() {
	if err := cliCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	addCommands()

	// Setup command flags
	cliCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config", "",
		"config file (default is $HOME/.config/plantd/plant.yaml)",
	)
	cliCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	viper.BindPFlag("verbose", cliCmd.PersistentFlags().Lookup("verbose"))
	viper.SetDefault("verbose", false)
}

func addCommands() {
	// cliCmd.AddCommand(jobCmd)
	cliCmd.AddCommand(stateCmd)

	// Miscellaneous commands
	cliCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	err := core.LoadConfig("client", &config)
	if err != nil {
		fmt.Errorf("Fatal error reading config file: %s \n", err)
	}

	endpoint = config.Server.Endpoint
}
