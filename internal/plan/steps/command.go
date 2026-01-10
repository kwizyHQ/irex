package steps

import (
	"os"
	"os/exec"

	"github.com/kwizyHQ/irex/internal/plan"
)

type CommandStep struct {
	Args                []string
	RunWhen             *bool
	DescriptionOverride string
}

func (c *CommandStep) ID() string {
	return "run:command"
}

func (c *CommandStep) Name() string {
	return "Command Execution"
}

func (c *CommandStep) Description() string {
	if c.DescriptionOverride != "" {
		return c.DescriptionOverride
	}
	return "Executes a command in the target directory."
}

func (c *CommandStep) Run(ctx *plan.PlanContext) error {
	if c.RunWhen != nil && !*c.RunWhen {
		return nil
	}
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Dir = ctx.TargetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
