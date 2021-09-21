package main

func init() {
	// Make roots version option only emit the version. This is used in Actions.
	// The template looks weird on purpose. Leaving as a single line causes the
	// output to append an extra character.
	rootCmd.Version = "0.0.1"
	rootCmd.SetVersionTemplate(
		`{{printf "%s" .Version}}`)
}
