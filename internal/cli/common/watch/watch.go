package watch

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	nodets "github.com/kwizyHQ/irex/internal/engines/node-ts"
	"github.com/kwizyHQ/irex/internal/ir"
	"github.com/kwizyHQ/irex/internal/plan"
	steps "github.com/kwizyHQ/irex/internal/plan/steps"
	"github.com/kwizyHQ/irex/internal/tempdir"
	"github.com/kwizyHQ/irex/internal/watcher"
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
			dir := tempdir.Get()
			planCtx := plan.PlanContext{
				TargetDir:         ".",
				IR:                &ir.IRBundle{},
				TmpDir:            dir,
				CompiledTemplates: make(plan.CompiledTemplates),
				RenderSession:     &plan.RenderSession{},
			}

			watchPlan := &plan.Plan{
				ID:   "watch",
				Name: "Watch server",
				Steps: []plan.Step{
					&steps.LoadIR{IRPath: "irex.hcl"},
					&steps.PlanSelectorStep{
						PlansMap: map[string]func(ctx *plan.PlanContext) *plan.Plan{
							"node-ts": nodets.NodeTSWatchPlan,
						},
						DeferLoadingKey: func(psCtx *plan.PlanContext) string {
							return psCtx.IR.Config.Runtime.Name
						},
					},
				},
			}

			// 🔹 1. Initial execution
			slog.Debug("Initial build")
			if err := watchPlan.Execute(&planCtx); err != nil {
				slog.Error("initial build failed", "err", err)
				os.Exit(1)
			}

			// 🔹 2. Start watcher
			mgr := watcher.NewManager(
				[]string{
					"irex.hcl",
					planCtx.IR.Config.Paths.Specifications + "/**/*",
					planCtx.IR.Config.Paths.Templates + "/**/*",
				},
				300*time.Millisecond,
				func(ctx context.Context, events []watcher.Event) error {
					slog.Debug("Change detected, rebuilding", "events", len(events))
					for _, ev := range events {
						slog.Debug(" - "+ev.Path, "type", ev.Type)
					}
					// IMPORTANT: reuse same PlanContext
					// Later you can diff IR / runtime here
					return watchPlan.Execute(&planCtx)
				},
				false,
			)

			go func() {
				if err := mgr.Run(ctx); err != nil {
					slog.Error("watcher stopped", "err", err)
				}
			}()

			slog.Info("Watcher running")
			<-ctx.Done()
			// Cleanup on exit
			slog.Info("Shutting down watcher")
			dir.Delete()
		},
	}

	return cmd
}
