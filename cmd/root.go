package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	root := &cobra.Command{
		Use:   "nuon-ext-starship",
		Short: "Starship prompt module for Nuon CLI context",
		Long: `Display the current Nuon org, install, and environment in your Starship prompt.

Commands:
  prompt    Output the formatted prompt string (called by Starship)
  init      Print the starship.toml snippet to add
  refresh   Force-refresh the name cache`,
	}

	root.AddCommand(
		promptCmd(),
		initCmd(),
		refreshCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
