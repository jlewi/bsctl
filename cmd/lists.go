package cmd

import (
	"fmt"
	"github.com/jlewi/bsctl/pkg/application"
	"github.com/jlewi/bsctl/pkg/lists"
	"github.com/jlewi/bsctl/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

// NewDumpListCmd create a command to dump the list to a YAML file
func NewDumpListCmd() *cobra.Command {
	var (
		listURI    string
		outputPath string
	)

	// TODO(jeremy): We should update apply to support the image resource.
	cmd := &cobra.Command{
		Use:   "dumpList -o <file to output to> -u <uri of the list>",
		Short: "Merge the list of accounts in the starter pack into the specified YAML file.",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				defer app.Shutdown()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				version.Log()

				if err := app.SetupRegistry(); err != nil {
					return err
				}

				client, err := app.GetXRPCClient()
				if err != nil {
					return err
				}

				return lists.MergeListToFile(client, outputPath, listURI)
			}()
			if err != nil {
				fmt.Printf("Error running apply;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&listURI, "uri", "u", "", "List URI (required)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path of the output YAML file (required)")

	// Mark flags as required
	cmd.MarkFlagRequired("uri")
	cmd.MarkFlagRequired("output")

	return cmd
}
