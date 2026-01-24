package validate

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// ValidateSchema performs semantic checks on a ModelsSpec and returns diagnostics for all issues found.
func ValidateSchema(spec *symbols.ModelsSpec) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}

	if spec == nil || spec.ModelsBlock == nil {
		reporter.Error("Missing required 'models' block.", zeroRange, "irex.input.required", "models")
		return reporter.All()
	}
	block := spec.ModelsBlock
	if len(block.Models) == 0 {
		reporter.Warn("No models defined in 'models' block.", zeroRange, "irex.input.recommended", "models")
	}
	modelNames := map[string]struct{}{}
	for _, model := range block.Models {
		if model.Name == "" {
			reporter.Error("Model 'name' is required.", zeroRange, "irex.input.required", "models.name")
			continue
		}
		if _, exists := modelNames[model.Name]; exists {
			reporter.Error("Duplicate model name: "+model.Name, zeroRange, "irex.input.duplicate", "models.name")
		} else {
			modelNames[model.Name] = struct{}{}
		}
		if len(model.Fields) == 0 {
			reporter.Error("Model '"+model.Name+"' must have at least one field.", zeroRange, "irex.input.required", "models.fields")
		}
		for _, field := range model.Fields {
			checkModelFieldSemantics(field, model.Name, reporter, zeroRange)
		}
		// Relations block checks (optional)
		if model.Relations != nil {
			for _, rel := range model.Relations.HasMany {
				if rel.Name == "" {
					reporter.Error("hasMany relation in model '"+model.Name+"' missing name.", zeroRange, "irex.input.required", "models.relation.hasMany.name")
				}
				if rel.Ref == "" {
					reporter.Error("hasMany relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "irex.input.required", "models.relation.hasMany.ref")
				}
			}
			for _, rel := range model.Relations.BelongsTo {
				if rel.Name == "" {
					reporter.Error("belongsTo relation in model '"+model.Name+"' missing name.", zeroRange, "irex.input.required", "models.relation.belongsTo.name")
				}
				if rel.Ref == "" {
					reporter.Error("belongsTo relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "irex.input.required", "models.relation.belongsTo.ref")
				}
			}
			for _, rel := range model.Relations.ManyToMany {
				if rel.Name == "" {
					reporter.Error("manyToMany relation in model '"+model.Name+"' missing name.", zeroRange, "irex.input.required", "models.relation.manyToMany.name")
				}
				if rel.Ref == "" {
					reporter.Error("manyToMany relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "irex.input.required", "models.relation.manyToMany.ref")
				}
			}
		}
		// Config block checks (optional)
		if model.Config != nil {
			// if model.Config.IDStrategy == "" {
			// 	reporter.Info("Model '"+model.Name+"' config: idStrategy is not set (using default).", zeroRange, "irex.input.recommended", "models.config.idStrategy")
			// }
			if model.Config.DB != nil {
				// Example: warn if both mongo and mysql are empty
				if model.Config.DB.Mongo == (symbols.MongoDBConfig{}) && model.Config.DB.Mysql == (symbols.MySqlDBConfig{}) {
					reporter.Warn("Model '"+model.Name+"' config.db: both mongo and mysql configs are empty.", zeroRange, "irex.input.recommended", "models.config.db")
				}
			}
		}
	}
	return reporter.All()
}

func checkModelFieldSemantics(field symbols.ModelField, modelName string, reporter *diagnostics.Reporter, rng diagnostics.Range) {
	if field.Name == "" {
		reporter.Error("Field in model '"+modelName+"' missing name.", rng, "irex.input.required", "models.fields.name")
	}
	if field.Type == "" && len(field.Fields) == 0 {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"' must have a type or nested fields.", rng, "irex.input.required", "models.fields.type_or_nested")
	}
	if field.MinLength != nil && field.MaxLength != nil && *field.MinLength > *field.MaxLength {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"': minlength > maxlength.", rng, "irex.input.invalid", "models.fields.length")
	}
	if field.Min != nil && field.Max != nil && *field.Min > *field.Max {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"': min > max.", rng, "irex.input.invalid", "models.fields.range")
	}
	// if field.Unique && field.DB != nil && field.DB.Mongo != nil && !field.DB.Mongo.Unique {
	//      reporter.Warn("Field '"+field.Name+"' in model '"+modelName+"' is unique but mongo db config does not set unique.", rng, "irex.input.mismatch", "models.fields.db.mongo.unique")
	// }
	for _, nested := range field.Fields {
		checkModelFieldSemantics(nested, modelName, reporter, rng)
	}
}
