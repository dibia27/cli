package db

import (
	"errors"

	"github.com/spf13/cobra"
)

var size, dbName, networkID, software, softwareVersion, firewallID string
var replicas, numSnapshots int
var publicIP bool

// DBCmd is the root command for the db subcommand
var DBCmd = &cobra.Command{
	Use:     "db",
	Aliases: []string{"database"},
	Short:   "Manage Civo Database ",
	Long:    `Create, update, delete, and list Civo Database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	DBCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&size, "size", "s", "g3.medium", "the size of nodes to create.")
	createCmd.Flags().StringVarP(&networkID, "network", "n", "default", "the name of the network to use in the creation")
	// TODO: Should we call it database type? So the short hand would be "t"
	createCmd.Flags().StringVarP(&software, "software", "w", "mysql", "the name of the software to use in the creation")
	createCmd.Flags().StringVarP(&softwareVersion, "software-version", "v", "", "the version of the software to use in the creation")
	createCmd.Flags().IntVarP(&replicas, "replicas", "r", 0, "the number of replicas in addition to the primary node")
	createCmd.Flags().IntVarP(&numSnapshots, "num-snapshots", "", 1, "the number of snapshots to keep")
	createCmd.Flags().BoolVarP(&publicIP, "public-ip", "p", true, "whether to assign a public IP to the database")
	createCmd.Flags().StringVarP(&firewallID, "firewall", "", "", "the firewall to use for the database")

	DBCmd.AddCommand(getCommand)
	DBCmd.AddCommand(listCmd)
}
