package cmd

import (
	"context"
	"fmt"

	"github.com/nuonco/nuon-ext-starship/internal/cache"
	"github.com/nuonco/nuon-ext-starship/internal/config"
	"github.com/nuonco/nuon-ext-starship/internal/resolve"
	"github.com/spf13/cobra"
)

func refreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Force-refresh the cached org and install names",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRefresh()
		},
	}
}

func runRefresh() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if cfg.APIToken == "" {
		return fmt.Errorf("no API token configured")
	}
	if cfg.OrgID == "" {
		return fmt.Errorf("no org selected")
	}

	ctx := context.Background()

	entry := &cache.Entry{
		OrgID:  cfg.OrgID,
		APIURL: cfg.APIURL,
	}

	orgName, err := resolve.OrgName(ctx, cfg.APIURL, cfg.APIToken, cfg.OrgID)
	if err != nil {
		return fmt.Errorf("resolving org name: %w", err)
	}
	entry.OrgName = orgName
	fmt.Printf("org: %s (%s)\n", orgName, cfg.OrgID)

	if cfg.InstallID != "" {
		installName, err := resolve.InstallName(ctx, cfg.APIURL, cfg.APIToken, cfg.OrgID, cfg.InstallID)
		if err != nil {
			return fmt.Errorf("resolving install name: %w", err)
		}
		entry.InstallID = cfg.InstallID
		entry.InstallName = installName
		fmt.Printf("install: %s (%s)\n", installName, cfg.InstallID)
	}

	if err := cache.Write(entry); err != nil {
		return fmt.Errorf("writing cache: %w", err)
	}

	fmt.Println("cache updated ✓")
	return nil
}
