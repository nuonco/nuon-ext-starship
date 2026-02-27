package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Print the starship.toml snippet to add to your config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(starshipSnippet)
		},
	}
}

const starshipSnippet = `# Add to your ~/.config/starship.toml:

[custom.nuon]
command = "nuon starship prompt"
detect_files = ["runner.toml"]
format = "[$symbol($output )]($style)"
symbol = "☁️ "
style = "bold purple"
`
