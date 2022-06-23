package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var format string

var objectStoreCredentialExportCmd = &cobra.Command{
	Use:     "export",
	Aliases: []string{"export-credentials"},
	Short:   "Export the credentials for your Object Store.",
	Example: "civo objectstore credential export --access-key=ACCESS_KEY --format=FORMAT (We support env and s3cfg)",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		var key string
		if accessKey != "" {
			if format == "" {
				utility.Error("You must provide a format to export to. See --help for more information.")
				os.Exit(1)
			}
			key = accessKey
		}
		if format != "" {
			if accessKey == "" {
				utility.Error("You must provide an access key. See --help for more information.")
				os.Exit(1)
			}
		}
		if key == "" {
			utility.Error("You must provide an access key and the format to export to. See --help for more information.")
			os.Exit(1)
		}

		objectStore, err := client.FindObjectStore(key)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if format == "env" {
			fmt.Printf("# Tip: You can redirect output with (>> ~/.zshrc) to add these to Zsh's startup automatically\n")
			fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", objectStore.AccessKeyID)
			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", objectStore.SecretAccessKey)
			fmt.Printf("export AWS_DEFAULT_REGION=%s\n", client.Region)
			// TODO: change this to objectstore.region.civo.com before GA
			fmt.Printf("export AWS_HOST=https://objectstorage.%s.civo.io\n", strings.ToLower(client.Region))
		} else if format == "s3cfg" {
			fmt.Printf("# Tip: You can redirect output with (>> ~/.s3cfg) to automatically configure s3cmd\n")
			fmt.Printf("[default]\n")
			fmt.Printf("access_key = %s\n", objectStore.AccessKeyID)
			fmt.Printf("access_token = %s\n", objectStore.SecretAccessKey)
			fmt.Printf("bucket_location = %s\n", client.Region)
			fmt.Printf("host_base = objectstorage.%s.civo.io\n", strings.ToLower(client.Region))
		} else {
			utility.Error("You must provide a valid format to export to. Supported formats are env and s3cfg. See --help for more information.")
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		}
	},
}
