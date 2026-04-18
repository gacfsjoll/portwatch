package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/user/portwatch/internal/acknowledge"
)

func runAcknowledge(args []string, storePath string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: portwatch acknowledge <port>")
	}
	port, err := parsePort(args[0])
	if err != nil {
		return err
	}
	s := acknowledge.NewStore(storePath)
	if err := s.Load(); err != nil {
		return fmt.Errorf("load acks: %w", err)
	}
	if err := s.Acknowledge(port); err != nil {
		return fmt.Errorf("acknowledge: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Port %d acknowledged.\n", port)
	return nil
}

func runAcknowledgeList(storePath string) error {
	s := acknowledge.NewStore(storePath)
	if err := s.Load(); err != nil {
		return fmt.Errorf("load acks: %w", err)
	}
	list := s.List()
	if len(list) == 0 {
		fmt.Fprintln(os.Stdout, "No acknowledged ports.")
		return nil
	}
	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
	fmt.Fprintln(os.Stdout, "Acknowledged ports:")
	for _, p := range list {
		fmt.Fprintf(os.Stdout, "  %d\n", p)
	}
	return nil
}

func runAcknowledgeRevoke(args []string, storePath string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: portwatch acknowledge revoke <port>")
	}
	port, err := parsePort(args[0])
	if err != nil {
		return err
	}
	s := acknowledge.NewStore(storePath)
	if err := s.Load(); err != nil {
		return fmt.Errorf("load acks: %w", err)
	}
	if err := s.Revoke(port); err != nil {
		return fmt.Errorf("revoke: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Acknowledgement for port %d revoked.\n", port)
	return nil
}

func parsePort(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("invalid port %q: %w", s, err)
	}
	return uint16(n), nil
}
