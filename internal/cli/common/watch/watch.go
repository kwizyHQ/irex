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

var (
	comCtx = context.Background()
)

func Run() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch mode (placeholder)",
		Run: func(cmd *cobra.Command, args []string) {
			// 🔴 DO NOT use signal.NotifyContext (Cobra exits early on Windows)
			ctx, cancel := context.WithCancel(comCtx)

			sigCh := make(chan os.Signal, 1)
			signal.Ignore(os.Interrupt)
			signal.Notify(sigCh, os.Interrupt)

			dir := tempdir.Get()
			planCtx := plan.PlanContext{
				TargetDir:         ".",
				IR:                &ir.IRBundle{},
				TmpDir:            dir,
				CompiledTemplates: make(plan.CompiledTemplates),
				RenderSession:     &plan.RenderSession{},
				WatchRegistry:     plan.NewWatchRegistry(),
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

			// 🔹 Initial build
			slog.Debug("Initial build")
			if err := watchPlan.Execute(&planCtx); err != nil {
				slog.Error("initial build failed", "err", err)
				os.Exit(1)
			}

			mgr := watcher.NewManager(
				[]string{
					"irex.hcl",
					planCtx.IR.Config.Paths.Specifications + "/**/*",
					planCtx.IR.Config.Paths.Templates + "/**/*",
				},
				300*time.Millisecond,
				func(ctx context.Context, events []watcher.Event) error {
					if ctx.Err() != nil {
						return nil // prevent rebuild during shutdown
					}
					slog.Debug("Change detected, rebuilding", "events", len(events))
					return watchPlan.Execute(&planCtx)
				},
				false,
			)

			go func() {
				if err := mgr.Run(ctx); err != nil {
					slog.Error("watcher stopped", "err", err)
				}
			}()

			// 🔥 CENTRALIZED SHUTDOWN (THIS IS THE FIX)
			go func() {
				<-sigCh
				slog.Warn("Ctrl+C received, shutting down")

				// 1️⃣ stop watcher loop
				cancel()

				// 2️⃣ stop running dev processes
				planCtx.WatchRegistry.Shutdown()

				// 4️⃣ cleanup temp dir
				dir.Delete()

				os.Exit(0)
			}()

			slog.Info("Watcher running")
			select {} // block forever, we exit explicitly
		},
	}
	cmd.SetContext(comCtx)
	return cmd
}
