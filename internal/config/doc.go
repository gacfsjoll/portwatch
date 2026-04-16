// Package config provides loading, validation, and default values for
// portwatch configuration files.
//
// Configuration is expressed as YAML and supports the following top-level
// sections:
//
//	# portwatch.yaml
//	scan:
//	  port_start: 1
//	  port_end:   65535
//	  interval:   30s
//	alert:
//	  backend: log      # "log" or "stdout"
//	  log_file: ""      # optional; defaults to stderr when empty
//
// # Loading configuration
//
// Use [Load] to read a configuration file from disk. If the file does not
// exist, [Load] returns an error wrapping [os.ErrNotExist].
//
// Use [Default] to obtain a configuration with sensible out-of-the-box
// values without requiring a file on disk. This is useful for tests or
// single-binary deployments that embed their settings.
//
// # Validation
//
// Both [Load] and [Default] return a validated [Config]. Validation ensures
// that port_start is less than or equal to port_end, that both values fall
// within the valid port range (1–65535), and that the chosen alert backend
// is recognised.
package config
