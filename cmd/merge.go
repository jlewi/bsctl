package cmd

import (
	"fmt"
	"github.com/jlewi/bsctl/pkg/api/v1alpha1"
	"github.com/jlewi/bsctl/pkg/application"
	"github.com/jlewi/bsctl/pkg/lists"
	"github.com/jlewi/bsctl/pkg/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

// NewMergeCmd merge two lists
func NewMergeCmd() *cobra.Command {

	var source string
	var destFile string
	var sourceFilter string
	cmd := &cobra.Command{
		Use:   "merge --source=<resource.yaml> --dest=<resourceDir> --source-filter=all|members|non-members",
		Short: "Merge the two lists",
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

				filter := lists.IncludeAll
				switch sourceFilter {
				case "all":
					filter = lists.IncludeAll
				case "members":
					filter = lists.IncludeMembers
				case "nonmembers":
					filter = lists.IncludeNonMembers
				default:
					return fmt.Errorf("invalid source filter: %v; value should be all, members or nonmembers", sourceFilter)
				}

				srcB, err := os.ReadFile(source)
				if err != nil {
					return errors.Wrapf(err, "Failed to read source file; %v", source)
				}

				src := &v1alpha1.AccountList{}
				if err := yaml.Unmarshal(srcB, src); err != nil {
					return errors.Wrapf(err, "Failed to decode source file; %v", source)
				}

				destB, err := os.ReadFile(destFile)
				if err != nil {
					return errors.Wrapf(err, "Failed to read dest file; %v", destFile)
				}

				dest := &v1alpha1.AccountList{}
				if err := yaml.Unmarshal(destB, dest); err != nil {
					return errors.Wrapf(err, "Failed to decode dest file; %v", destFile)
				}

				lists.MergeFollowLists(dest, *src, filter)

				outB, err := yaml.Marshal(dest)
				if err != nil {
					return errors.Wrapf(err, "Failed to encode merged AccountList")
				}

				if err := os.WriteFile(destFile, outB, 0644); err != nil {
					return errors.Wrapf(err, "Failed to write merged AccountList to file; %v", destFile)
				}
				return nil
			}()
			if err != nil {
				fmt.Printf("Error running apply;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&source, "source", "s", "", "The source list")
	cmd.Flags().StringVarP(&destFile, "dest", "d", "", "The destination list")
	cmd.Flags().StringVarP(&sourceFilter, "source-filter", "", "all", "The source filter; value can be all, members, or nonmembers")

	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("dest")

	return cmd
}
