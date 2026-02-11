package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "",
				Usage:   "The JSON file to provide as input.",
			},
		},
		Action: WithFile,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		return
	}
}
