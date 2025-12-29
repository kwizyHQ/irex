package assemble

import (
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

// prepareServiceIR converts a symbols.Service into an ir.IRService and
// registers it in the bundle. It also collects middlewares/policies at service level
// into the IRRoutes and IROperations later via apply blocks.
func prepareServiceIR(ctx *shared.BuildContext, svc symbols.Service, parent string) error {
	if ctx == nil {
		return nil
	}
	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.Services == nil {
		ctx.IR.Services = make(ir.IRServices)
	}

	name := svc.Name
	kind := ir.ServiceCustom
	if svc.Model != "" {
		kind = ir.ServiceModel
	}

	rs := ir.IRService{
		Name:   name,
		Kind:   kind,
		Model:  svc.Model,
		Parent: parent,
		Expose: svc.Expose,
	}

	ctx.IR.Services[name] = rs

	return nil
}
