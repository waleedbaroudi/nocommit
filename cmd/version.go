package cmd

// Version is overridden at build time using -ldflags.
// Defaults to "dev" for local builds.
var Version = "dev"

func init() {
	RootCmd.Version = Version
}
