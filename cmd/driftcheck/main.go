package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/driftcheck/internal/compose"
	"github.com/driftcheck/internal/docker"
	"github.com/driftcheck/internal/drift"
	"github.com/driftcheck/internal/output"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		composeFile = flag.String("file", "docker-compose.yml", "path to compose file")
		formatFlag  = flag.String("format", "text", "output format: text or json")
		exitCode    = flag.Bool("exit-code", false, "exit with code 1 if drift is found")
	)
	flag.Parse()

	project, err := compose.ParseFile(*composeFile)
	if err != nil {
		return fmt.Errorf("parsing compose file: %w", err)
	}

	client, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	defer client.Close()

	containers, err := client.ListContainers(context.Background())
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	serviceMap := docker.BuildServiceMap(containers)
	report := drift.Detect(project, serviceMap)

	fmt := output.NewFormatter(os.Stdout, output.Format(*formatFlag))
	if err := fmt.Write(report); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	if *exitCode && report.Summary().Total > 0 {
		os.Exit(1)
	}
	return nil
}
