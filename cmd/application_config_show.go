package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var appConfigShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show application config",
	Example: "civo app config show APP_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		app, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, config := range app.Config {
			ow.StartLine()
			ow.AppendDataWithLabel("name", config.Name, "Name")
			ow.AppendDataWithLabel("value", config.Value, "Value")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
