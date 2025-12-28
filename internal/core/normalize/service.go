package normalize

import (
	"fmt"
	"reflect"

	"github.com/kwizyHQ/irex/internal/core/symbols"
)

func SafeSet(parentField, childField reflect.Value) {
	if !childField.CanSet() {
		return
	}

	// Check if we are dealing with a pointer
	if parentField.Kind() == reflect.Ptr && !parentField.IsNil() {
		// Create a BRAND NEW memory address of the same type
		newValue := reflect.New(parentField.Type().Elem())

		// Copy the value from the parent address to the new address
		newValue.Elem().Set(parentField.Elem())

		// Point the child to our new, independent memory address
		childField.Set(newValue)
	} else {
		// For non-pointers (int, string, etc.), Set() copies by value automatically
		childField.Set(parentField)
	}
}

// MergeDefaults merges parent values into child. parent is parentDefaults, and child is childDefaults.
// Child must be a pointer so it can be modified.
func MergeDefaults(parent any, child any) {
	// Get the underlying values
	pVal := reflect.ValueOf(parent)
	cVal := reflect.ValueOf(child)

	// If they are pointers, get the element they point to
	if pVal.Kind() == reflect.Ptr {
		pVal = pVal.Elem()
	}
	if cVal.Kind() == reflect.Ptr {
		cVal = cVal.Elem()
	}

	// Safety check: must be structs of the same type
	if pVal.Kind() != reflect.Struct || cVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < cVal.NumField(); i++ {
		childField := cVal.Field(i)
		parentField := pVal.Field(i)

		// If child field is empty (Zero Value), take from parent
		if childField.IsZero() {
			SafeSet(parentField, childField)
		}
	}
}

// MergeFromDefaults takes a 'source' (defaults) and a 'target' (to be updated).
// It searches for fields in target that match the name and type of fields in source.
func MergeFromDefaults(source any, target any) {
	srcVal := reflect.ValueOf(source)
	tgtVal := reflect.ValueOf(target)

	// We must have a pointer to the target to modify it
	if tgtVal.Kind() != reflect.Ptr {
		fmt.Println("Error: Target must be a pointer")
		return
	}
	tgtVal = tgtVal.Elem()

	// If source is a pointer, get the underlying struct
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	// Iterate through every field in the Source (the defaults)
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		fieldName := srcVal.Type().Field(i).Name
		// skip if source field is zero value
		if srcField.IsZero() {
			continue
		}

		// Try to find a field with the same name in the Target
		targetField := tgtVal.FieldByName(fieldName)

		// Verification:
		// 1. Does the field exist in target?
		// 2. Is it settable (exported)?
		// 3. Is it currently empty (Zero Value)?
		// 4. Do the types match?
		// print info about target and src field like name value type etc
		if targetField.IsValid() && targetField.CanSet() && targetField.IsZero() {
			SafeSet(srcField, targetField)
		}
	}
}

func NormalizeServiceAST(def *symbols.ServiceDefinition) {
	if def == nil || def.Services == nil {
		return
	}

	if def.RateLimits != nil && def.RateLimits.Defaults != nil {
		for i := range def.RateLimits.Presets {
			// Use the address of the slice element (not the loop copy)
			MergeFromDefaults(def.RateLimits.Defaults, &def.RateLimits.Presets[i])
		}
	}

	// now we need to walk the service tree and apply defaults recursively
	var walk func(svcs []symbols.Service, parentDefaults *symbols.ServiceDefaults)
	walk = func(svcs []symbols.Service, parentDefaults *symbols.ServiceDefaults) {
		for i := range svcs {
			svc := &svcs[i]
			// if svc.Defaults is nil, initialize it
			if svc.Defaults == nil {
				svc.Defaults = &symbols.ServiceDefaults{}
			}
			// Step 1: Merge parent defaults into current service
			MergeDefaults(parentDefaults, svc.Defaults)
			// Step 2: Apply explicit service overrides
			MergeFromDefaults(svc.Defaults, svc)

			// Step 3: Recurse into child services
			if len(svc.Services) > 0 {
				walk(svc.Services, svc.Defaults)
			}
		}
	}
	walk(def.Services.Services, def.Services.Defaults)

}
