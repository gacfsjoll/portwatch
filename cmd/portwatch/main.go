// Command portwatch is a lightweight daemon that monitors open TCP ports
// and alerts when unexpected changes are detected.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	cfg", "portwatch.yaml", "path to configuration file")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}

	notifier, err := alert.FromConfig(cfg.Alert.Backend, cfg.Alert.LogFile)
	if err != nil {
		log.Fatalf("portwatch: alert setup: %v", err)
	}

	sc := scanner.NewScanner(cfg.Scan.PortStart, cfg.Scan.PortEnd)
	m := monitor.New(sc, notifier, cfg.Scan.Interval)

	go func() {
		if err := m.Run(); err != nil {
			log.Fatalf("portwatch: monitor: %v", err)
		}
	}()

	log.Printf("portwatch: monitoring ports %d-%d every %v",
		cfg.Scan.PortStart, cfg.Scan.PortEnd, cfg.Scan.Interval)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("portwatch: shutting down")
	m.Stop()
}

// loadConfig attempts to load the file at path; if the file does not exist
// it falls back to built-in defaults so portwatch works without a config file.
func loadConfig(path string) (config.Config, error) {
	cfg, err := config.Load(path)
	if os.IsNotExist(err) {
		log.Printf("portwatch: config file %q not found, using defaults", path)
		return config.Default(), nil
	}
	return cfg, err
}
