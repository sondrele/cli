package validatecmd

import (
	"fmt"
	"github.com/urfave/cli/v2"

	"github.com/nais/cli/pkg/validate"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "validate",
		Usage:           "Validate nais.yaml configuration",
		ArgsUsage:       "nais.yaml [naiser.yaml...]",
		UsageText:       "nais validate nais.yaml [naiser.yaml...]",
		HideHelpCommand: true,
		Action: func(context *cli.Context) error {
			if context.Args().Len() == 0 {
				return fmt.Errorf("no config files provided")
			}

			return validate.NaisConfig(context.Args().Slice())
		},
	}
}
