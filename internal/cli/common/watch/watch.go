package watch

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	nodets "github.com/kwizyHQ/irex/internal/engines/node-ts"
	"github.com/kwizyHQ/irex/internal/ir"
	"github.com/kwizyHQ/irex/internal/plan"
	. "github.com/kwizyHQ/irex/internal/plan/steps"
	"github.com/spf13/cobra"
)

func Run() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch mode (placeholder)",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()
			/**
			* Watcher functionality to be implemented
			* LoadIR (from root of the project irex.hcl)
			* Compile templates
			* Render output
			* Watch for file changes and re-compile as needed
			 */
			// let's first implement a normal build command before watch (using plan)

			planCtx := plan.PlanContext{
				TargetDir: ".",
				IR:        &ir.IRBundle{},
			}

			watchPlan := &plan.Plan{
				ID:   "watch",
				Name: "Watch server",
				Steps: []plan.Step{
					&LoadIR{IRPath: "irex.hcl"},
					&PlanSelectorStep{
						PlansMap: map[string]func(ctx *plan.PlanContext) *plan.Plan{
							"node-ts": nodets.NodeTSWatchPlan,
						},
						GetByKey: func(psCtx *plan.PlanContext) string {
							return psCtx.IR.Config.Runtime.Name
						},
					},
				},
			}

			slog.Info("Starting the watcher.")
			err := watchPlan.Execute(&planCtx)

			if err != nil {
				slog.Error(err.Error())
			}

			<-ctx.Done()
		},
	}

	return cmd
}
