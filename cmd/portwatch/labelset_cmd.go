package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/portwatch/internal/labelset"
)

// runLabelSetValidate is a small diagnostic sub-command that validates and
// prints a normalised label set supplied on the command line.
//
// Usage: portwatch labels env=prod owner=alice tier=backend
func runLabelSetValidate(args []string) error {
	fs := flag.NewFlagSet("labels", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: portwatch labels [key=value ...]")
		fmt.Fprintln(os.Stderr, "  Validates and prints normalised labels.")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	pairs := fs.Args()
	if len(pairs) == 0 {
		fs.Usage()
		return fmt.Errorf("no labels provided")
	}

	ls, err := labelset.New(pairs...)
	if err != nil {
		return fmt.Errorf("invalid label set: %w", err)
	}

	fmt.Fprintf(os.Stdout, "labels: %s\n", ls.String())

	all := ls.All()
	keys := make([]string, 0, len(all))
	for k := range all {
		keys = append(keys, k)
	}
	// Print individual labels for human readability.
	for _, k := range sortedKeys(all) {
		fmt.Fprintf(os.Stdout, "  %-20s = %s\n", k, all[k])
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// simple insertion sort – label sets are tiny
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && strings.Compare(keys[j-1], keys[j]) > 0; j-- {
			keys[j-1], keys[j] = keys[j], keys[j-1]
		}
	}
	return keys
}
