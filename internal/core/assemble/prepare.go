package assemble

import (
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

// PrepareIR prepares an IRBundle from parsed ASTs available in ctx.
// It currently maps HTTP settings, services and operations into the IR.
func ProjectIR(ctx *shared.BuildContext) error {
	if ctx == nil {
		return nil
	}

	// delegate to specialized preparers
	if err := preparePoliciesIR(ctx); err != nil {
		return err
	}
	if err := prepareRateLimitsIR(ctx); err != nil {
		return err
	}
	if err := prepareConfigIR(ctx); err != nil {
		return err
	}
	if err := prepareModelsIR(ctx); err != nil {
		return err
	}

	if ctx.IR == nil {
		bundle := ir.IRBundle{}

		ctx.IR = &bundle
		return nil
	}

	s := ctx.ServicesAST.Services

	// HTTP config (only HTTP mapping here)
	ctx.IR.Http = ir.IRHttpConfig{
		BasePath:         s.BasePath,
		Cors:             s.Cors,
		AllowedOrigins:   s.AllowedOrigins,
		AllowedMethods:   s.AllowedMethods,
		AllowedHeaders:   s.AllowedHeaders,
		ExposeHeaders:    s.ExposeHeaders,
		AllowCredentials: s.AllowCredentials,
		MaxAge:           s.MaxAge,
		CacheControl:     s.CacheControl,
	}

	// now let's do walking services and operations and prepare routes, operations, services, middlewares, etc.
	// we walk through services recursively to capture nested services and so execute the logic in one pass.
	type ServiceWalkContext struct {
		Services      *[]symbols.Service
		Operations    *[]symbols.Operation
		ParentService string
	}

	var walk func(walkCtx *ServiceWalkContext)
	walk = func(walkCtx *ServiceWalkContext) {
		if walkCtx.Services == nil {
			return
		}
		for i := range *walkCtx.Services {
			svc := &(*walkCtx.Services)[i]
			// process service call prepareServiceIR
			prepareServiceIR(ctx, svc, walkCtx.ParentService)
			// infer operations from service in case of model-based service
			prepareInferredOperationsIR(ctx, svc, walkCtx.ParentService)
			walk(&ServiceWalkContext{
				Services:      &svc.Services,
				Operations:    &svc.Operations,
				ParentService: svc.Name,
			})
		}
		// process operations at this level
		if walkCtx.Operations != nil {
			for i := range *walkCtx.Operations {
				op := &(*walkCtx.Operations)[i]
				prepareOperationIR(ctx, op, walkCtx.ParentService)
			}
		}
	}
	walk(&ServiceWalkContext{
		Services:   &s.Services,   // top-level services
		Operations: &s.Operations, // top-level operations
	})

	// ctx.IR = &bundle

	return nil
}
