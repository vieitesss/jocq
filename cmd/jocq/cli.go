package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
	"github.com/vieitesss/jocq/internal/buffer"
	"github.com/vieitesss/jocq/internal/ingest"
	"github.com/vieitesss/jocq/internal/tui"
)

func RunTUI(inputFile string) error {
	opts := []tea.ProgramOption{tea.WithAltScreen()}

	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create the data store
	data := buffer.NewData()

	// Create the scanner
	scanner := ingest.NewScanner(f, data)
	// We'll start by reading everything first
	// This will end up being a goroutine
	err = scanner.Scan()
	if err != nil {
		return err
	}

	// Start the TUI
	app := tui.NewApp(data)
	_, err = tea.NewProgram(app, opts...).Run()
	if err != nil {
		return err
	}

	return nil
}

func WithFile(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.String("file")

	if inputFile == "" {
		return fmt.Errorf("You must provide a JSON file with the `--file` flag.")
	}

	return RunTUI(inputFile)
}
