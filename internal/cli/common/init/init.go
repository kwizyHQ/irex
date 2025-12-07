package initcmd

import (
	"github.com/spf13/cobra"

	nodeTsCmd "github.com/kwizyHQ/irex/internal/engines/node-ts/bootstrap"
)

// NewInitCmd returns a Cobra command for project initialization
func Run() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a dev project (interactive)",
	}

	cmd.AddCommand(nodeTsCmd.Run())

	return cmd
}
