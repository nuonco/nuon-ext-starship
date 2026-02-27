package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/nuonco/nuon-ext-starship/internal/cache"
	"github.com/nuonco/nuon-ext-starship/internal/config"
	"github.com/nuonco/nuon-ext-starship/internal/resolve"
	"github.com/spf13/cobra"
)

func promptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompt",
		Short: "Output the formatted prompt string for Starship",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrompt()
		},
	}
}

func runPrompt() error {
	cfg, err := config.Load()
	if err != nil {
		return nil // silent fail — don't break the prompt
	}

	if cfg.OrgID == "" {
		return nil // nothing selected
	}

	env := cfg.EnvName()

	// Check cache.
	cached, _ := cache.Read()

	orgName := ""
	installName := ""
	needsRefresh := false

	if cached != nil && cached.OrgID == cfg.OrgID && cached.APIURL == cfg.APIURL {
		orgName = cached.OrgName
		if cfg.InstallID != "" && cached.InstallID == cfg.InstallID {
			installName = cached.InstallName
		} else if cfg.InstallID != "" {
			needsRefresh = true
		}
	} else {
		needsRefresh = true
	}

	if needsRefresh && cfg.APIToken != "" {
		ctx := context.Background()

		if orgName == "" {
			name, err := resolve.OrgName(ctx, cfg.APIURL, cfg.APIToken, cfg.OrgID)
			if err == nil {
				orgName = name
			}
		}

		if cfg.InstallID != "" && installName == "" {
			name, err := resolve.InstallName(ctx, cfg.APIURL, cfg.APIToken, cfg.OrgID, cfg.InstallID)
			if err == nil {
				installName = name
			}
		}

		// Update cache.
		entry := &cache.Entry{
			OrgID:       cfg.OrgID,
			OrgName:     orgName,
			InstallID:   cfg.InstallID,
			InstallName: installName,
			APIURL:      cfg.APIURL,
		}
		_ = cache.Write(entry)
	}

	// Fallback to truncated IDs if names couldn't be resolved.
	if orgName == "" {
		orgName = truncateID(cfg.OrgID)
	}
	if cfg.InstallID != "" && installName == "" {
		installName = truncateID(cfg.InstallID)
	}

	// Build output.
	var parts []string
	parts = append(parts, env)
	parts = append(parts, "🏢 "+orgName)
	if installName != "" {
		parts = append(parts, "📦 "+installName)
	}

	fmt.Print(strings.Join(parts, " "))
	return nil
}

// truncateID returns the last 8 characters of an ID as a fallback display.
func truncateID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return "…" + id[len(id)-8:]
}
