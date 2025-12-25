package semantic

import (
	"github.com/kwizyHQ/irex/internal/core/symbols"
	"github.com/kwizyHQ/irex/internal/diagnostics"
)

// CheckSchemaSemantics performs semantic checks on a ModelsSpec and returns diagnostics for all issues found.
func CheckSchemaSemantics(spec *symbols.ModelsSpec) []Diagnostic {
	reporter := diagnostics.NewReporter()
	zeroRange := diagnostics.Range{}
	source := "schema"

	if spec == nil || spec.ModelsBlock == nil {
		reporter.Error("Missing required 'models' block.", zeroRange, "schema.models.missing", source)
		return reporter.All()
	}
	block := spec.ModelsBlock
	if len(block.Models) == 0 {
		reporter.Warn("No models defined in 'models' block.", zeroRange, "schema.models.empty", source)
	}
	modelNames := map[string]struct{}{}
	for _, model := range block.Models {
		if model.Name == "" {
			reporter.Error("Model 'name' is required.", zeroRange, "schema.model.name.required", source)
			continue
		}
		if _, exists := modelNames[model.Name]; exists {
			reporter.Error("Duplicate model name: "+model.Name, zeroRange, "schema.model.name.duplicate", source)
		} else {
			modelNames[model.Name] = struct{}{}
		}
		if len(model.Fields) == 0 {
			reporter.Error("Model '"+model.Name+"' must have at least one field.", zeroRange, "schema.model.fields.required", source)
		}
		for _, field := range model.Fields {
			checkModelFieldSemantics(field, model.Name, reporter, source, zeroRange)
		}
		// Relations block checks (optional)
		if model.Relations != nil {
			for _, rel := range model.Relations.HasMany {
				if rel.Name == "" {
					reporter.Error("hasMany relation in model '"+model.Name+"' missing name.", zeroRange, "schema.model.relation.hasMany.name.required", source)
				}
				if rel.Ref == "" {
					reporter.Error("hasMany relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "schema.model.relation.hasMany.ref.required", source)
				}
			}
			for _, rel := range model.Relations.BelongsTo {
				if rel.Name == "" {
					reporter.Error("belongsTo relation in model '"+model.Name+"' missing name.", zeroRange, "schema.model.relation.belongsTo.name.required", source)
				}
				if rel.Ref == "" {
					reporter.Error("belongsTo relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "schema.model.relation.belongsTo.ref.required", source)
				}
			}
			for _, rel := range model.Relations.ManyToMany {
				if rel.Name == "" {
					reporter.Error("manyToMany relation in model '"+model.Name+"' missing name.", zeroRange, "schema.model.relation.manyToMany.name.required", source)
				}
				if rel.Ref == "" {
					reporter.Error("manyToMany relation '"+rel.Name+"' in model '"+model.Name+"' missing ref.", zeroRange, "schema.model.relation.manyToMany.ref.required", source)
				}
			}
		}
		// Config block checks (optional)
		if model.Config != nil {
			// if model.Config.IDStrategy == "" {
			// 	reporter.Info("Model '"+model.Name+"' config: idStrategy is not set (using default).", zeroRange, "schema.model.config.idStrategy.default", source)
			// }
			if model.Config.DB != nil {
				// Example: warn if both mongo and mysql are empty
				if model.Config.DB.Mongo == (symbols.MongoDBConfig{}) && model.Config.DB.Mysql == (symbols.MySqlDBConfig{}) {
					reporter.Warn("Model '"+model.Name+"' config.db: both mongo and mysql configs are empty.", zeroRange, "schema.model.config.db.empty", source)
				}
			}
		}
	}
	return reporter.All()
}

func checkModelFieldSemantics(field symbols.ModelField, modelName string, reporter *diagnostics.Reporter, source string, rng diagnostics.Range) {
	if field.Name == "" {
		reporter.Error("Field in model '"+modelName+"' missing name.", rng, "schema.model.field.name.required", source)
	}
	if field.Type == "" && len(field.Fields) == 0 {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"' must have a type or nested fields.", rng, "schema.model.field.type_or_nested.required", source)
	}
	if field.MinLength != nil && field.MaxLength != nil && *field.MinLength > *field.MaxLength {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"': minlength > maxlength.", rng, "schema.model.field.length.invalid", source)
	}
	if field.Min != nil && field.Max != nil && *field.Min > *field.Max {
		reporter.Error("Field '"+field.Name+"' in model '"+modelName+"': min > max.", rng, "schema.model.field.range.invalid", source)
	}
	// if field.Unique && field.DB != nil && field.DB.Mongo != nil && !field.DB.Mongo.Unique {
	// 	reporter.Warn("Field '"+field.Name+"' in model '"+modelName+"' is unique but mongo db config does not set unique.", rng, "schema.model.field.db.mongo.unique.mismatch", source)
	// }
	for _, nested := range field.Fields {
		checkModelFieldSemantics(nested, modelName, reporter, source, rng)
	}
}
