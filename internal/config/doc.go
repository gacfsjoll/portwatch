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
// Use [Load] to read a file from disk, or [Default] to obtain a configuration
// with sensible out-of-the-box values.
package config
