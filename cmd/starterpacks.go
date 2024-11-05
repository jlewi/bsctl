package cmd

import (
	"fmt"
	"github.com/jlewi/bsctl/pkg/application"
	"github.com/jlewi/bsctl/pkg/lists"
	"github.com/jlewi/bsctl/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

// NewDumpStarterPack create a command to dump the starter pack to a YAML file
func NewDumpStarterPack() *cobra.Command {
	var (
		accountHandle   string
		starterPackName string
		outputPath      string
	)

	// TODO(jeremy): We should update apply to support the image resource.
	cmd := &cobra.Command{
		Use:   "dumpStarterPack -o <file to output to> -h <handle of account that owns it> -n <name of the starter pack>",
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

				return lists.MergeStarterPackToFile(client, outputPath, accountHandle, starterPackName)
			}()
			if err != nil {
				fmt.Printf("Error running apply;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&accountHandle, "handle", "a", "", "Handle of the account (required)")
	cmd.Flags().StringVarP(&starterPackName, "name", "n", "", "Name of the starter pack (required)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path of the output YAML file (required)")

	// Mark flags as required
	cmd.MarkFlagRequired("account")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("output")

	return cmd
}
