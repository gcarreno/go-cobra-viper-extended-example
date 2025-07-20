package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/gcarreno/go-cobra-viper-extended-example/config"
)

// Always define constants early, in order to make changes quick and depend
// on syntax completion to avoid errors
const (
	cFlagLogLevel      = "l"
	cFlagLogLevelLong  = "log-level"
	cFlagLogLevelUsage = "log level: info, warn, error, debug"

	cFlagAdminEmail      = "a"
	cFlagAdminEmailLong  = "admin-email"
	cFlagAdminEmailUsage = "site admin email"

	cFlagWebAddressLong  = "web-address"
	cFlagWebAddressUsage = "web server address"

	cFlagWebPortLong  = "web-port"
	cFlagWebPortUsage = "web server port: [1024, 65535]"

	cFlagAPIAddressLong  = "api-address"
	cFlagAPIAddressUsage = "api server address"

	cFlagAPIPortLong  = "api-port"
	cFlagAPIPortUsage = "api server port: [1, 65535]"
)

var (
	logLevel   string
	adminEmail string
	apiAddress string
	apiPort    int32
	webAddress string
	webPort    int32

	logLevels = map[string]bool{
		"info":  true,
		"warn":  true,
		"error": true,
		"debug": true,
	}

	logLevelsCompletions = []cobra.Completion{
		cobra.CompletionWithDesc("info", "debug level info"),
		cobra.CompletionWithDesc("warn", "debug level warn"),
		cobra.CompletionWithDesc("error", "debug level error"),
		cobra.CompletionWithDesc("debug", "debug level info"),
	}

	// serveCmd represents the serve command
	serveCmd = &cobra.Command{
		// The name of this command
		Use: "serve",

		// Define some aliases
		Aliases: []string{"s"},

		Short: "Starts the web server",
		// 	Long: `A longer description that spans multiple lines and likely contains examples
		// and usage of using your command. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,

		// Some examples
		Example: `  # All defaults
  mysite serve`,

		// We use PreRunE in order to validate that the flags contain allowed values
		PreRunE: servePreRunE,

		// Run: ,
		// We use RunE instead of Run, so that we can trigger the usage to show when
		// there's an error, this following how cobra works when a flag error is triggered.
		RunE: serveRunE,
	}
)

// Add the serve command to the command chain and register any flags for this command
func init() {
	rootCmd.AddCommand(serveCmd)

	// Obtain the defaults
	config := cfg.DefaultConfig()

	// Flag log level
	serveCmd.Flags().StringVarP(&logLevel, cFlagLogLevelLong, cFlagLogLevel, config.LogLevel, cFlagLogLevelUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperLogLevel, serveCmd.Flags().Lookup(cFlagLogLevelLong))
	// Register the completion values for flag "log-level"
	err := serveCmd.RegisterFlagCompletionFunc(
		cFlagLogLevelLong,
		cobra.FixedCompletions(logLevelsCompletions, cobra.ShellCompDirectiveNoFileComp),
	)
	if err != nil {
		// Use rootCmd.Print and others to print to the same place cobra does
		serveCmd.PrintErrf("error registering flag completion function: %v\n", err)
		os.Exit(1)
	}

	// Flag admin email
	serveCmd.Flags().StringVarP(&adminEmail, cFlagAdminEmailLong, cFlagAdminEmail, config.AdminEmail, cFlagAdminEmailUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperAdminEmail, serveCmd.Flags().Lookup(cFlagAdminEmailLong))

	// Flag api address
	serveCmd.Flags().StringVar(&apiAddress, cFlagAPIAddressLong, config.API.Address, cFlagAPIAddressUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperAPIAddress, serveCmd.Flags().Lookup(cFlagAPIAddressLong))

	// Flag api port
	serveCmd.Flags().Int32Var(&apiPort, cFlagAPIPortLong, config.API.Port, cFlagAPIPortUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperAPIPort, serveCmd.Flags().Lookup(cFlagAPIPortLong))

	// Flag web address
	serveCmd.Flags().StringVar(&webAddress, cFlagWebAddressLong, config.Web.Address, cFlagWebAddressUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperWebAddress, serveCmd.Flags().Lookup(cFlagWebAddressLong))

	// Flag web port
	serveCmd.Flags().Int32Var(&webPort, cFlagWebPortLong, config.Web.Port, cFlagWebPortUsage)
	// Bind the flag to viper
	viper.BindPFlag(cfg.ViperWebPort, serveCmd.Flags().Lookup(cFlagWebPortLong))
}

// Check if anything is wrong with the flags before we run the main code of the command
func servePreRunE(cmd *cobra.Command, args []string) error {
	// Initialize viper's config
	if cfgFile == "" {
		// Look for default file
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.myapp")
	} else {
		// User gave us a config file
		viper.SetConfigFile(cfgFile)
	}

	// Load .env file
	err := godotenv.Load() // Automatically loads ".env"
	if err != nil {
		cmd.Println("No .env file found (that's okay)")
	}

	// Enable ENV binding
	viper.SetEnvPrefix("MYSITE")
	viper.AutomaticEnv()

	// Make env vars like MYSITE_API_ADDRESS map to "api.address"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read it
	if err := viper.ReadInConfig(); err == nil {
		cmd.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		cmd.Printf("No config file found, continuing with flags/env/defaults: %v\n", err)
	}

	// Get the defaults
	config = cfg.DefaultConfig()

	// Unmarshal it into the config structure
	if err := viper.Unmarshal(&config); err != nil {
		// This is an internal error, so we don't make cobra print the usage
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	if err := validateLogLevel(); err != nil {
		return err
	}

	if err := validateAPIPort(); err != nil {
		return err
	}

	if err := validateWebPort(); err != nil {
		return err
	}

	return nil
}

// The serve command main code
func serveRunE(cmd *cobra.Command, args []string) error {
	cmd.Printf("cfgFile: '%s'\n", cfgFile)
	cmd.Printf("Log Level: '%s'\n", config.LogLevel)
	cmd.Printf("Admin Email: '%s'\n", config.AdminEmail)
	cmd.Printf("API Address: '%s'\n", config.API.Address)
	cmd.Printf("API Port: '%d'\n", config.API.Port)
	cmd.Printf("Web Address: '%s'\n", config.Web.Address)
	cmd.Printf("Web Port: '%d'\n", config.Web.Port)

	return nil
}

func validateLogLevel() error {
	if !logLevels[strings.ToLower(config.LogLevel)] {
		return fmt.Errorf("log level must be one of: info, warn, error, debug (got: %s)", config.LogLevel)
	}

	return nil
}

// Validate the web port
func validateWebPort() error {
	// The web server should not running as root, so limit to [1024, 65535]
	if webPort < 1024 || webPort > 65535 {
		return fmt.Errorf("web port should be between [1024, 65535]")
	}

	return nil
}

// Validate the API port
func validateAPIPort() error {
	// The api is remote, so limit to [1, 65535]
	if apiPort < 1 || apiPort > 65535 {
		return fmt.Errorf("web port should be between [1, 65535]")
	}

	return nil
}
