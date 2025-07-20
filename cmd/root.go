package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	cfg "github.com/gcarreno/go-cobra-viper-extended-example/config"
)

// Always define constants early, in order to make changes quick and depend
// on syntax completion to avoid errors
const (
	cFlagConfigFile      = "c"
	cFlagConfigFileLong  = "config"
	cFlagConfigFileUsage = "config file (YAML, TOML, or JSON)"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	config  *cfg.Config
	rootCmd = &cobra.Command{
		// Name of the application
		Use: "mysite",

		// The version of this application
		Version: "0.0.1",

		Short: "My Web Site Server",
		// 	Long: `A longer description that spans multiple lines and likely contains
		// examples and usage of using your application. For example:

		// Cobra is a CLI library for Go that empowers applications.
		// This application is a tool to generate the needed files
		// to quickly create a Cobra application.`,

		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func (cmd *cobra.Command, args []string) {},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil { // This is important to track !!
		os.Exit(1)
	}
}

// Register any flags for the root command
func init() {
	// Register persistent flags that is present in any command
	// Config File flag
	rootCmd.PersistentFlags().StringVarP(&cfgFile, cFlagConfigFileLong, cFlagConfigFile, "", cFlagConfigFileUsage)

	// Register our own usage function
	rootCmd.SetUsageFunc(usage)
}

// This will print our own version of the usage
// IMPORTANT: This has to be checked with the default usage function to maintain
// the same result that is printed by cobra's default usage function and
// the default usage template
func usage(cmd *cobra.Command) error {
	// Use cmd.Print and others to print to the same place cobra does
	cmd.Print("\033[1mUSAGE\033[0m")
	if cmd.Runnable() {
		cmd.Printf("\n  %s", cmd.UseLine())
	}
	if cmd.HasAvailableSubCommands() {
		cmd.Printf("\n  %s [command]", cmd.CommandPath())
	}
	if len(cmd.Aliases) > 0 {
		cmd.Printf("\n\n\033[1mALIASES\033[0m\n")
		cmd.Printf("  %s", cmd.NameAndAliases())
	}
	if cmd.HasExample() {
		cmd.Printf("\n\n\033[1mEXAMPLES\033[0m\n")
		cmd.Printf("%s", cmd.Example)
	}
	if cmd.HasAvailableSubCommands() {
		cmds := cmd.Commands()
		if len(cmd.Groups()) == 0 {
			cmd.Printf("\n\n\033[1mAVAILABLE COMMANDS\033[0m")
			for _, subcmd := range cmds {
				if subcmd.IsAvailableCommand() || subcmd.Name() == "help" {
					cmd.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
				}
			}
		} else {
			for _, group := range cmd.Groups() {
				cmd.Printf("\n\n%s", group.Title)
				for _, subcmd := range cmds {
					if subcmd.GroupID == group.ID && (subcmd.IsAvailableCommand() || subcmd.Name() == "help") {
						cmd.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
			if !cmd.AllChildCommandsHaveGroup() {
				cmd.Printf("\n\n\033[1mADDITIONAL COMMANDS\033[0m")
				for _, subcmd := range cmds {
					if subcmd.GroupID == "" && (subcmd.IsAvailableCommand() || subcmd.Name() == "help") {
						cmd.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
		}
	}
	if cmd.HasAvailableLocalFlags() {
		cmd.Printf("\n\n\033[1mFLAGS\033[0m\n")
		cmd.Print(trimRightSpace(cmd.LocalFlags().FlagUsages()))
	}
	if cmd.HasAvailableInheritedFlags() {
		cmd.Printf("\n\n\033[1mGLOBAL FLAGS\033[0m\n")
		cmd.Print(trimRightSpace(cmd.InheritedFlags().FlagUsages()))
	}
	if cmd.HasHelpSubCommands() {
		cmd.Printf("\n\n\033[1mADDITIONAL HELP TOPICS\033[0m")
		for _, subcmd := range cmd.Commands() {
			if subcmd.IsAdditionalHelpTopicCommand() {
				cmd.Printf("\n  %s %s", rpad(subcmd.CommandPath(), subcmd.CommandPathPadding()), subcmd.Short)
			}
		}
	}

	if cmd.HasAvailableSubCommands() {
		cmd.Printf("\n\nUse \"%s [command] --help\" for more information about a command.", cmd.CommandPath())
	}
	cmd.Println()
	return nil
}

// Helper function to trim spaces at the right of a string
func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// Helper function to right pad a string with spaces
func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}
