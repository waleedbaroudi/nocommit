package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type FieldType string

const (
	StringT FieldType = "string"
	BoolT   FieldType = "bool"
)

type Config struct {
	Version       string `yaml:"version"`
	CaseSensitive bool   `yaml:"caseSensitive"`
}

type Field struct {
	Type    FieldType
	Default any
	Set     func(*Config, any)
	Get     func(Config) any
}

var ConfigSchema = map[string]Field{
	"version": {
		Type:    StringT,
		Default: "1",
		Set:     func(c *Config, v any) { c.Version = v.(string) },
		Get:     func(c Config) any { return c.Version },
	},
	"caseSensitive": {
		Type:    BoolT,
		Default: false,
		Set:     func(c *Config, v any) { c.CaseSensitive = v.(bool) },
		Get:     func(c Config) any { return c.CaseSensitive },
	},
}

func DefaultsConfig() Config {
	return Config{
		Version:       ConfigSchema["version"].Default.(string),
		CaseSensitive: ConfigSchema["caseSensitive"].Default.(bool),
	}
}

func ParseValue(t FieldType, s string) (any, error) {
	switch t {
	case StringT:
		return s, nil
	case BoolT:
		b, err := strconv.ParseBool(strings.TrimSpace(s))
		if err != nil {
			return nil, fmt.Errorf("expected boolean true/false")
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unsupported type %q", t)
	}
}

func AllowedKeysCSV() string {
	ks := make([]string, 0, len(ConfigSchema))
	for k := range ConfigSchema {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return strings.Join(ks, ", ")
}

func loadConfig(path string) (Config, error) {
	var out Config
	b, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}

func SaveConfig(path string, cfg Config) error {
	b, err := yaml.Marshal(cfg) // struct order => version first
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}
