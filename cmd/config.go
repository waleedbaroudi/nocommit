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

		home, err := os.UserHomeDir()
		if err != nil {
			fail("Get home dir", err, true)
			return
		}
		configFile := filepath.Join(home, ".nocommit", "config.yaml")

		data, err := os.ReadFile(configFile)
		if err != nil {
			fail("Read config.yaml", err, true)
			return
		}

		cfg := make(map[string]interface{})
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fail("Parse config.yaml", err, true)
			return
		}

		// If only key provided → print value
		writer := cmd.OutOrStdout()
		if len(args) == 1 {
			if val, ok := cfg[key]; ok {
				fmt.Fprintf(writer, "%s: %v\n", key, val)
			} else {
				fmt.Fprintf(writer, "%s not set\n", key)
			}
			return
		}

		// Otherwise update and write back
		value := args[1]
		cfg[key] = value // TODO: disallow new keys?

		out, err := yaml.Marshal(cfg)
		if err != nil {
			fail("Encode config.yaml", err, true)
			return
		}

		if err := os.WriteFile(configFile, out, 0600); err != nil {
			fail("Write config.yaml", err, true)
			return
		}
		fmt.Fprintf(writer, "✅ Updated %s = %s\n", key, value)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fail("Get home dir", err, true)
			return
		}
		configFile := filepath.Join(home, ".nocommit", "config.yaml")

		data, err := os.ReadFile(configFile)
		if err != nil {
			fail("Read config.yaml", err, true)
			return
		}

		w := cmd.OutOrStdout()
		fmt.Fprintln(w, string(data))
	},
}

func init() {
	configCmd.AddCommand(configListCmd) // add "list" under "config"
	RootCmd.AddCommand(configCmd)       // add "config" under root
}
