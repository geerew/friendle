package cmd

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// rootCmd is the root of the Friendle CLI
var rootCmd = &cobra.Command{
	Use:   "friendle",
	Short: "Friendle",
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// init configures global viper settings shared by subcommands
func init() {
	viper.SetEnvPrefix("FRIENDLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Execute is the entry point for the CLI
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// questionPassword asks a question, waits for hidden text input and returns the answer
func questionPassword(question string) string {
	c := color.New(color.Bold, color.FgGreen)
	c.Printf(">> %s: ", question)

	answerBytes, _ := term.ReadPassword(int(os.Stdin.Fd()))
	answer := string(answerBytes)
	answer = strings.TrimSpace(answer)

	c.Println()

	return answer
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// errorMessage prints an error message
func errorMessage(message string, a ...any) {
	c := color.New(color.FgRed)
	c.Printf(message+"\n", a...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// successMessage prints a success message
func successMessage(message string, a ...any) {
	c := color.New(color.FgYellow)
	c.Printf(message+"\n", a...)
}
