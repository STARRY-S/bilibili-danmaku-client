package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/STARRY-S/bilibili-danmaku-client/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func initializeFlagsConfig(cmd *cobra.Command, cfg config.Provider) {
	if cmd.Parent() != nil {
		initializeFlagsConfig(cmd.Parent(), cfg)
	}

	flags := cmd.Flags()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		setValueFromFlag(flags, f.Name, cfg, "")
	})
}

func setValueFromFlag(flags *pflag.FlagSet, key string, cfg config.Provider, targetKey string) {
	key = strings.TrimSpace(key)
	if flags.Lookup(key) != nil || flags.Changed(key) {
		f := flags.Lookup(key)
		configKey := key
		if targetKey != "" {
			configKey = targetKey
		}
		// Gotta love this API.
		switch f.Value.Type() {
		case "bool":
			bv, _ := flags.GetBool(key)
			cfg.Set(configKey, bv)
		case "string":
			cfg.Set(configKey, f.Value.String())
		case "stringSlice":
			sv, _ := flags.GetStringSlice(key)
			cfg.Set(configKey, sv)
		case "int":
			iv, _ := flags.GetInt(key)
			cfg.Set(configKey, iv)
		default:
			panic(fmt.Sprintf("update switch with %s", f.Value.Type()))
		}
	}
}

func getConfigJson(cfg config.Provider) string {
	c := cfg.Get("")
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
