package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/gcarreno/go-cobra-viper-extended-example/config"
)

// Always define constants early, in order to make changes quick and depend
// on syntax completion to avoid errors
const (
	cFlagType      = "t"
	cFlagTypeLong  = "config-type"
	cFlagTypeUsage = "config type for the config file: toml, json, yaml, yml"
	cConfigName    = "config"
	cConfigType    = "toml"
)

var (
	configType  string
	configTypes = map[string]bool{
		"toml": true,
		"json": true,
		"yaml": true,
		"yml":  true,
	}
	configTypesCompletions = []cobra.Completion{
		cobra.CompletionWithDesc("toml", "toml format"),
		cobra.CompletionWithDesc("json", "JSON format"),
		cobra.CompletionWithDesc("yaml", "YAML format"),
		cobra.CompletionWithDesc("yml", "YAML format"),
	}

	// initCmd represents the init command
	initCmd = &cobra.Command{
		// The name of this command
		Use: "init",

		// Define some aliases
		Aliases: []string{"i"},

		Short: "Initializes and writes a config file",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,

		// Some examples
		Example: `  # All defaults
  mysite init

  # With a custom filename
  mysite init --config ~/.config.toml

  # With a different config type
  mysite init --config-type yml`,

		// We use PreRunE in order to validate that the flags contain allowed values
		// values
		PreRunE: runtInitPreRunE,

		// Run: ,
		// We use RunE instead of Run, so that we can trigger the usage to show when
		// there's an error, this following how cobra works when a flag error is triggered.
		RunE: runInitE,
	}
)

// Add the init command to the command chain and register any flags for this command
func init() {
	rootCmd.AddCommand(initCmd)

	// Flag "config-type"
	initCmd.Flags().StringVarP(&configType, cFlagTypeLong, cFlagType, cConfigType, cFlagTypeUsage)
	// Register the completion values for flag "config-type"
	err := initCmd.RegisterFlagCompletionFunc(
		cFlagTypeLong,
		cobra.FixedCompletions(configTypesCompletions, cobra.ShellCompDirectiveNoFileComp),
	)
	if err != nil {
		// Use rootCmd.Print and others to print to the same place cobra does
		initCmd.PrintErrf("error registering flag completion function: %v\n", err)
		os.Exit(1)
	}
}

// Check if anything is wrong with the flags before we run the main code of the command
// Good place to check if our config type flag contains a valid option
func runtInitPreRunE(cmd *cobra.Command, args []string) error {
	// We could just return the error from validateConfigType, but I'm allowing for more
	// validate functions if needed
	if err := validateConfigType(); err != nil {
		return err
	}

	return nil
}

// The init command main code
func runInitE(cmd *cobra.Command, args []string) error {
	// Setup the config filename
	if cfgFile == "" {
		viper.SetConfigName(cConfigName)
		viper.SetConfigType(configType)
		cfgFile = fmt.Sprintf("%s.%s", cConfigName, configType)
	} else {
		ext := filepath.Ext(cfgFile)
		if ext == "" {
			// We need to force an extension or else viper will give an error and
			// I don't want to include a config type flag on the serve command, or
			// even globally, for that matter
			switch configType {
			case "toml", "json", "yaml", "yml":
				cfgFile = fmt.Sprintf("%s.%s", cfgFile, configType)
			default:
				cfgFile = fmt.Sprintf("%s.%s", cfgFile, cConfigType)
			}
		}
		viper.SetConfigName(cfgFile)
	}

	// Load the defaults into viper
	cfg.SetDefaultsToViper()

	// Fill in our local config structure
	// This is only needed to print the JSON content down below
	// In production, this would not be necessary
	if err := viper.Unmarshal(&config); err != nil {
		// This is an internal error, so we don't make cobra print the usage
		cmd.PrintErrf("error unmarshaling config: %v\n", err)
		os.Exit(1)
	}

	if err := viper.SafeWriteConfigAs(cfgFile); err != nil {
		// This is an internal error, so we don't make cobra print the usage
		cmd.PrintErrf("error writing config: %v\n", err)
		os.Exit(1)
	}

	// Dump final config
	cmd.Println("Final Config:")
	jsonDump, _ := json.MarshalIndent(config, "", "  ")
	cmd.Println(string(jsonDump))

	return nil
}

// Validates if the config type flag contains a valid option
func validateConfigType() error {
	if !configTypes[strings.ToLower(configType)] {
		return fmt.Errorf("config type must be one of: toml, json, yaml, yml (got: %s)", configType)
	}

	return nil
}
