package main

import (
	"fmt"
	"github.com/jpbede/ratiocheck/internal/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	app := &cli.App{
		Name:        "ratio-check",
		Usage:       "Image to content ratio check",
		Description: "Microservice to check image to content ration of HTML pages",
		Version:     fmt.Sprintf("%s-%s", version, commit),
		Commands: []*cli.Command{
			commands.Listen(),
			commands.Check(),
		},
	}

	// run app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
