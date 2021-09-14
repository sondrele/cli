package cmd

import (
	"fmt"
	"github.com/nais/nais-d/cmd/helpers"
	"github.com/nais/nais-d/pkg/consts"
	"github.com/nais/nais-d/pkg/secrets"
	"github.com/spf13/cobra"
	"os"
)

var getCmd = &cobra.Command{
	Use:   "get [args] [flags]",
	Short: "The specified config format will get the secret and generated to 'current' location",
	Run: func(cmd *cobra.Command, args []string) {
		secretName, err := helpers.GetString(cmd, SecretNameFlag, args[0], true)
		if err != nil {
			fmt.Printf("getting %s: %s", SecretNameFlag, err)
			os.Exit(1)
		}

		team, err := helpers.GetString(cmd, TeamFlag, args[1], true)
		if err != nil {
			fmt.Printf("getting %s: %s", SecretNameFlag, err)
			os.Exit(1)
		}

		configType, err := helpers.GetString(cmd, ConfigFlag, "", false)
		if err != nil {
			fmt.Printf("getting %s: %s", ConfigFlag, err)
			os.Exit(1)
		}

		if configType != consts.ENV && configType != consts.ALL && configType != consts.KCAT {
			fmt.Printf("valid args: %s | %s | %s", consts.ENV, consts.KCAT, consts.ALL)
			os.Exit(1)
		}

		dest, err := helpers.GetString(cmd, DestFlag, "", false)
		if err != nil {
			fmt.Printf("getting %s: %s", DestFlag, err)
			os.Exit(1)
		}

		dest, err = helpers.DefaultDestination(dest)
		if err != nil {
			fmt.Printf("an error %s", err)
			os.Exit(1)
		}
		secrets.ExtractAndGenerateConfig(configType, dest, secretName, team)
	},
}