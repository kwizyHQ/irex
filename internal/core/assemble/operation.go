package assemble

import (
	"fmt"
	"path"
	"strings"

	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/ir"
)

// joinPath joins base and seg ensuring single slashes
func joinPath(base, seg string) string {
	if base == "" {
		base = "/"
	}
	if seg == "" || seg == "/" {
		return path.Clean(base)
	}
	return path.Clean(base + "/" + strings.TrimPrefix(seg, "/"))
}

// prepareInferredOperationsIR infers basic CRUD operations for model-based services.
// For now it will not create a full set, only a minimal example for 'list' and 'read'.
func prepareInferredOperationsIR(ctx *shared.BuildContext, svc *symbols.Service, parentService string) error {
	if svc.Model == "" {
		return nil
	}

	// joinPath joins base and seg ensuring single slashes
	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.Operations == nil {
		ctx.IR.Operations = make(ir.IROperations)
	}
	// determine requested CRUD operations (service-level overrides defaults)
	var crudOps []string
	if len(svc.CrudOperations) > 0 {
		crudOps = svc.CrudOperations
	} else if svc.Defaults != nil && len(svc.Defaults.CrudOperations) > 0 {
		crudOps = svc.Defaults.CrudOperations
	}

	// helper to check presence (case-insensitive)
	has := func(name string) bool {
		for _, v := range crudOps {
			if v == "*" {
				return true
			}
			if strings.EqualFold(v, name) {
				return true
			}
		}
		return false
	}

	// default path base
	basePath := svc.Path
	if basePath == "" {
		basePath = "/"
	}

	// generate CREATE
	if has("CREATE") {
		name := fmt.Sprintf("%s.create", svc.Name)
		path := joinPath(basePath, "/")
		op := ir.IROperation{
			Name:    name,
			Service: svc.Name,
			Method:  "POST",
			Path:    path,
			Kind:    ir.OperationKindData,
			Data: &ir.DataOperationMeta{
				Action:        ir.DataCreate,
				Target:        "single",
				ReturnsEntity: true,
			},
		}
		ctx.IR.Operations[name] = op
		_ = prepareRouteIR(ctx, name, "POST", path, svc.Name, nil)
	}

	// generate READ
	if has("READ") {
		name := fmt.Sprintf("%s.read", svc.Name)
		path := joinPath(basePath, ":id")
		op := ir.IROperation{
			Name:    name,
			Service: svc.Name,
			Method:  "GET",
			Path:    path,
			Kind:    ir.OperationKindData,
			Data: &ir.DataOperationMeta{
				Action:        ir.DataRead,
				Target:        "single",
				ReturnsEntity: true,
			},
		}
		ctx.IR.Operations[name] = op
		_ = prepareRouteIR(ctx, name, "GET", path, svc.Name, nil)
	}

	// generate UPDATE
	if has("UPDATE") {
		name := fmt.Sprintf("%s.update", svc.Name)
		path := joinPath(basePath, ":id")
		op := ir.IROperation{
			Name:    name,
			Service: svc.Name,
			Method:  "PATCH",
			Path:    path,
			Kind:    ir.OperationKindData,
			Data: &ir.DataOperationMeta{
				Action: ir.DataUpdate,
				Target: "single",
			},
		}
		ctx.IR.Operations[name] = op
		_ = prepareRouteIR(ctx, name, "PATCH", path, svc.Name, nil)
	}

	// generate DELETE
	if has("DELETE") {
		name := fmt.Sprintf("%s.delete", svc.Name)
		path := joinPath(basePath, ":id")
		op := ir.IROperation{
			Name:    name,
			Service: svc.Name,
			Method:  "DELETE",
			Path:    path,
			Kind:    ir.OperationKindData,
			Data: &ir.DataOperationMeta{
				Action: ir.DataDelete,
				Target: "single",
			},
		}
		ctx.IR.Operations[name] = op
		_ = prepareRouteIR(ctx, name, "DELETE", path, svc.Name, nil)
	}

	// generate LIST
	if has("LIST") {
		name := fmt.Sprintf("%s.list", svc.Name)
		path := joinPath(basePath, "/")
		paginated := false
		if svc.Pagination != nil {
			paginated = *svc.Pagination
		} else if svc.Defaults != nil && svc.Defaults.Pagination != nil {
			paginated = *svc.Defaults.Pagination
		}
		op := ir.IROperation{
			Name:    name,
			Service: svc.Name,
			Method:  "GET",
			Path:    path,
			Kind:    ir.OperationKindData,
			Data: &ir.DataOperationMeta{
				Action:      ir.DataList,
				Target:      "many",
				Paginated:   paginated,
				ReturnsList: true,
			},
		}
		ctx.IR.Operations[name] = op
		_ = prepareRouteIR(ctx, name, "GET", path, svc.Name, nil)
	}

	return nil
}

// prepareOperationIR converts a symbols.Operation into an ir.IROperation and
// registers routes by calling prepareRouteIR.
func prepareOperationIR(ctx *shared.BuildContext, op *symbols.Operation, serviceName string) error {
	if ctx == nil {
		return nil
	}
	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.Operations == nil {
		ctx.IR.Operations = make(ir.IROperations)
	}

	name := op.Name
	method := op.Method
	if method == "" {
		method = "GET"
	}
	path := op.Path
	if path == "" {
		path = "/"
	}

	irop := ir.IROperation{
		Name:        name,
		Service:     serviceName,
		Method:      method,
		Path:        path,
		Action:      op.Action,
		Description: op.Description,
		Kind:        ir.OperationKindCustom,
	}

	ctx.IR.Operations[name] = irop

	// create route for this operation
	if err := prepareRouteIR(ctx, name, method, path, serviceName, op); err != nil {
		return err
	}

	return nil
}
