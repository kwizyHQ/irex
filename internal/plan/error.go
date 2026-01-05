package plan

import "os"

type PlanError struct {
	StepID  string
	Message string
}

func (e *PlanError) Error() string {
	return "Plan Error in step " + e.StepID + ": " + e.Message
}

func (e *PlanError) ID() string {
	return "error:" + e.StepID
}

func (e *PlanError) Name() string {
	return "Plan Error"
}

func (e *PlanError) Description() string {
	return e.Message
}

func (e *PlanError) Run(ctx *PlanContext) error {
	// halt execution when this step is reached
	println(e.Error())
	os.Exit(1)
	return nil
}
