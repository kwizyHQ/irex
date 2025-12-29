package assemble

import (
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/ir"
)

func preparePoliciesIR(ctx *shared.BuildContext) error {
	if ctx == nil || ctx.ServicesAST == nil {
		return nil
	}

	s := ctx.ServicesAST

	// ensure maps
	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.RequestPolicies == nil {
		ctx.IR.RequestPolicies = make(ir.IRRequestPolicies)
	}
	if ctx.IR.ResourcePolicies == nil {
		ctx.IR.ResourcePolicies = make(ir.IRResourcePolicies)
	}

	// Presets
	if s.Policies != nil {
		for _, p := range s.Policies.Presets {
			name := p.Name
			effect := ir.PolicyAllow
			if p.Effect == "deny" {
				effect = ir.PolicyDeny
			}
			// decide scope: request vs resource. default to request
			scope := p.Scope
			if scope == "resource" {
				ctx.IR.ResourcePolicies[name] = ir.IRResourcePolicy{
					IRPolicyBase: ir.IRPolicyBase{
						Name:        name,
						Rule:        p.Rule,
						Effect:      effect,
						Description: p.Description,
					},
				}
			} else {
				ctx.IR.RequestPolicies[name] = ir.IRRequestPolicy{
					IRPolicyBase: ir.IRPolicyBase{
						Name:        name,
						Rule:        p.Rule,
						Effect:      effect,
						Description: p.Description,
					},
				}
			}
		}

		// Custom policies
		for _, c := range s.Policies.Customs {
			name := c.Name
			scope := c.Scope
			if scope == "resource" {
				ctx.IR.ResourcePolicies[name] = ir.IRResourcePolicy{
					IRPolicyBase: ir.IRPolicyBase{
						Name:        name,
						Rule:        "", // custom may not have rule here
						Effect:      ir.PolicyAllow,
						Description: c.Description,
					},
				}
			} else {
				ctx.IR.RequestPolicies[name] = ir.IRRequestPolicy{
					IRPolicyBase: ir.IRPolicyBase{
						Name:        name,
						Rule:        "",
						Effect:      ir.PolicyAllow,
						Description: c.Description,
					},
				}
			}
		}

		// Groups - groups reference policies by name, we'll not expand them here
		for _, g := range s.Policies.Groups {
			_ = g // groups are handled at application time (routes)
		}
	}

	return nil
}
