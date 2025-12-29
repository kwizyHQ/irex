package irbuilder

import (
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/ir"
)

// PrepareIR prepares an IRBundle from parsed ASTs available in ctx.
// It currently maps HTTP settings, services and operations into the IR.
func PrepareIR(ctx *shared.BuildContext) error {
	if ctx == nil {
		return nil
	}

	var bundle ir.IRBundle

	// Map HTTP / services-level config if present
	if ctx.ConfigAST != nil && ctx.ServicesAST == nil {
		// nothing: config currently doesn't contain HTTP settings
	}

	if ctx.ServicesAST != nil && ctx.ServicesAST.Services != nil {
		s := ctx.ServicesAST.Services

		// HTTP config
		bundle.Http = ir.IRHttpConfig{
			BasePath:         s.BasePath,
			Cors:             new(bool),
			AllowedOrigins:   s.AllowedOrigins,
			AllowedMethods:   s.AllowedMethods,
			AllowedHeaders:   s.AllowedHeaders,
			ExposeHeaders:    s.ExposeHeaders,
			AllowCredentials: new(bool),
			MaxAge:           new(int),
			CacheControl:     s.CacheControl,
		}
		if s.Cors != nil {
			bundle.Http.Cors = s.Cors
		}
		if s.AllowCredentials != nil {
			bundle.Http.AllowCredentials = s.AllowCredentials
		}
		if s.MaxAge != nil {
			bundle.Http.MaxAge = s.MaxAge
		}

		// Services
		bundle.Services = make(ir.IRServices)
		for _, svc := range s.Services {
			name := svc.Name
			bundle.Services[name] = ir.IRService{
				Name:  name,
				Model: svc.Model,
			}

			// Operations defined on the service
			for _, op := range svc.Operations {
				opname := op.Name
				if bundle.Operations == nil {
					bundle.Operations = make(ir.IROperations)
				}
				bundle.Operations[opname] = ir.IROperation{
					Name:        opname,
					Service:     svc.Name,
					Method:      op.Method,
					Path:        op.Path,
					Action:      op.Action,
					Description: op.Description,
				}
			}
		}

		// Top-level operations
		if s.Operations != nil {
			if bundle.Operations == nil {
				bundle.Operations = make(ir.IROperations)
			}
			for _, op := range s.Operations {
				opname := op.Name
				bundle.Operations[opname] = ir.IROperation{
					Name:        opname,
					Method:      op.Method,
					Path:        op.Path,
					Action:      op.Action,
					Description: op.Description,
				}
			}
		}
	}

	// Ensure maps are not nil to simplify downstream code
	if bundle.Services == nil {
		bundle.Services = make(ir.IRServices)
	}
	if bundle.Operations == nil {
		bundle.Operations = make(ir.IROperations)
	}
	if bundle.Routes == nil {
		bundle.Routes = make(ir.IRRoutes)
	}
	if bundle.Middlewares == nil {
		bundle.Middlewares = make(ir.IRMiddlewares)
	}
	if bundle.RequestPolicies == nil {
		bundle.RequestPolicies = make(ir.IRRequestPolicies)
	}
	if bundle.ResourcePolicies == nil {
		bundle.ResourcePolicies = make(ir.IRResourcePolicies)
	}
	if bundle.RateLimits == nil {
		bundle.RateLimits = make(ir.IRRateLimits)
	}

	ctx.IR = &bundle

	return nil
}
