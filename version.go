package main

// Version information - injected at build time via ldflags
var (
	// Version is the semantic version of the action
	Version = "dev"
	// Commit is the git commit hash
	Commit = "unknown"
)

// Action metadata
const (
	// ActionName is the full name of the GitHub Action
	ActionName = "appleboy/LLM-action"
	// ActionShortName is the short name used in User-Agent
	ActionShortName = "LLM-action"
)

// GetVersion returns the current version string
func GetVersion() string {
	return Version
}

// GetUserAgent returns the User-Agent string for HTTP requests
func GetUserAgent() string {
	return ActionShortName + "/" + Version
}
