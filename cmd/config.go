package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config [key] [value]",
	Short: "View or set configuration values",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		writer := cmd.OutOrStdout()

		home, err := os.UserHomeDir()
		if err != nil { fail("Get home dir", err, true); return }
		configFile := filepath.Join(home, ".nocommit", "config.yaml")

		cfg, err := loadConfig(configFile)
		if err != nil { fail("Read/parse config.yaml", err, true); return }

		// GET
		if len(args) == 1 {
			if field, ok := ConfigSchema[key]; ok {
				fmt.Fprintf(writer, "%s: %v\n", key, field.Get(cfg))
			} else {
				fmt.Fprintf(writer, "%s not set\n", key)
			}
			return
		}

		// SET via schema
		val := args[1]
		field, ok := ConfigSchema[key]
		if !ok {
			fmt.Fprintf(writer, "❌ Unknown key. Allowed: %s\n", AllowedKeysCSV())
			return
		}
		parsed, err := ParseValue(field.Type, val)
		if err != nil {
			fmt.Fprintf(writer, "❌ Invalid value for %s: %v (expected %s)\n", key, err, field.Type)
			return
		}
		field.Set(&cfg, parsed)

		if err := SaveConfig(configFile, cfg); err != nil {
			fail("Write config.yaml", err, true); return
		}
		fmt.Fprintf(writer, "✅ Updated %s = %v\n", key, field.Get(cfg))
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		writer := cmd.OutOrStdout()

		home, err := os.UserHomeDir()
		if err != nil { fail("Get home dir", err, true); return }
		configFile := filepath.Join(home, ".nocommit", "config.yaml")

		cfg, err := loadConfig(configFile)
		if err != nil { fail("Read/parse config.yaml", err, true); return }

		b, _ := yaml.Marshal(cfg) // preserves field order
		fmt.Fprintln(writer, string(b))
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	RootCmd.AddCommand(configCmd)
}
