package steps

import (
	"github.com/kwizyHQ/irex/internal/plan"
)

type WatchCommandStep struct {
	IDValue             string
	Args                []string
	RunWhen             *bool
	DescriptionOverride string
}

func (c *WatchCommandStep) ID() string {
	if c.IDValue != "" {
		return c.IDValue
	}
	if len(c.Args) > 0 {
		return "watch:" + c.Args[0]
	}
	return "watch:command"
}

func (c *WatchCommandStep) Name() string {
	return "Watch Command"
}

func (c *WatchCommandStep) Description() string {
	if c.DescriptionOverride != "" {
		return c.DescriptionOverride
	}
	return "Starts or restarts a long-running command and keeps it registered for future restarts."
}

func (c *WatchCommandStep) Run(ctx *plan.PlanContext) error {
	if c.RunWhen != nil && !*c.RunWhen {
		return nil
	}

	id := c.ID()
	return ctx.WatchRegistry.StartOrRestart(id, ctx.TargetDir, c.Args)
}
