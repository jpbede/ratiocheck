package commands

import (
	"errors"
	"fmt"
	"github.com/jpbede/ratiocheck/pkg/ratio"
	"github.com/urfave/cli/v2"
)

func Check() *cli.Command {
	return &cli.Command{
		Name:      "check",
		Aliases:   []string{"c"},
		ArgsUsage: "<URL to run the check on>",
		Usage:     "Run a one-off check for given url",
		Action:    runCheck,
	}
}

func runCheck(c *cli.Context) error {
	url := c.Args().First()
	if url == "" {
		return errors.New("missing url")
	}

	result, err := ratio.Get(c.Context, url)
	if err == nil {
		fmt.Printf("Ratio: %f %%\n", result.Ratio)
		fmt.Printf("Image Area: %f px2\n", result.ImageArea)
		fmt.Printf("Content Area: %f px2\n", result.ContentArea)
	}

	return err
}
