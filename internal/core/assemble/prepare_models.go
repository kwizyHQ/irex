package assemble

import (
	"github.com/kwizyHQ/irex/internal/core/shared"
	"github.com/kwizyHQ/irex/internal/ir"
)

// prepareModelsIR converts the parsed schema AST (models) into IRModels
// and stores them on ctx.IR.Models.
func prepareModelsIR(ctx *shared.BuildContext) error {
	if ctx == nil || ctx.SchemaAST == nil || ctx.SchemaAST.ModelsBlock == nil {
		return nil
	}

	if ctx.IR == nil {
		ctx.IR = &ir.IRBundle{}
	}
	if ctx.IR.Models == nil {
		ctx.IR.Models = make(ir.IRModels)
	}

	for _, m := range ctx.SchemaAST.ModelsBlock.Models {
		model := ir.IRModel{
			Name: m.Name,
		}

		// fields
		for _, f := range m.Fields {
			mf := ir.IRModelField{
				Name:        f.Name,
				Type:        f.Type,
				Required:    f.Required,
				Unique:      f.Unique,
				Trim:        f.Trim,
				MinLength:   f.MinLength,
				MaxLength:   f.MaxLength,
				Min:         f.Min,
				Max:         f.Max,
				Default:     f.Default,
				Match:       f.Match,
				Message:     f.Message,
				Visibility:  f.Visibility,
				Description: f.Description,
			}
			// nested fields
			if len(f.Fields) > 0 {
				for _, nf := range f.Fields {
					nmf := ir.IRModelField{
						Name:        nf.Name,
						Type:        nf.Type,
						Required:    nf.Required,
						Unique:      nf.Unique,
						Trim:        nf.Trim,
						MinLength:   nf.MinLength,
						MaxLength:   nf.MaxLength,
						Min:         nf.Min,
						Max:         nf.Max,
						Default:     nf.Default,
						Match:       nf.Match,
						Message:     nf.Message,
						Visibility:  nf.Visibility,
						Description: nf.Description,
					}
					mf.Fields = append(mf.Fields, nmf)
				}
			}

			// DB configs
			if f.DB != nil {
				dbc := &ir.IRModelFieldDBConfig{}
				if f.DB.Mongo != nil {
					dbc.Mongo = &ir.IRMongoFieldConfig{
						Index:  f.DB.Mongo.Index,
						Unique: f.DB.Mongo.Unique,
					}
					if f.DB.Mongo.Collation != nil {
						dbc.Mongo.Collation = &ir.IRMongoFieldCollation{
							Locale:          f.DB.Mongo.Collation.Locale,
							CaseLevel:       f.DB.Mongo.Collation.CaseLevel,
							CaseFirst:       f.DB.Mongo.Collation.CaseFirst,
							Strength:        f.DB.Mongo.Collation.Strength,
							NumericOrdering: f.DB.Mongo.Collation.NumericOrdering,
							Alternate:       f.DB.Mongo.Collation.Alternate,
							MaxVariable:     f.DB.Mongo.Collation.MaxVariable,
							Backwards:       f.DB.Mongo.Collation.Backwards,
						}
					}
				}
				if f.DB.Mysql != nil {
					dbc.Mysql = &ir.IRMySQLFieldConfig{
						Index:   f.DB.Mysql.Index,
						Unique:  f.DB.Mysql.Unique,
						Collate: f.DB.Mysql.Collate,
					}
				}
				mf.DB = dbc
			}

			model.Fields = append(model.Fields, mf)
		}

		// config
		if m.Config != nil {
			cfg := &ir.IRModelConfig{
				Timestamps:  m.Config.Timestamps,
				Table:       m.Config.Table,
				Strict:      m.Config.Strict,
				IDStrategy:  m.Config.IDStrategy,
				Description: m.Config.Description,
			}
			// indexes
			for _, idx := range m.Config.Indexes {
				cfg.Indexes = append(cfg.Indexes, ir.IRModelIndex{
					Name:   idx.Name,
					Fields: idx.Fields,
					Unique: idx.Unique,
				})
			}
			// DB
			if m.Config.DB != nil {
				dbcfg := &ir.IRModelConfigDB{}
				// mongo
				dbcfg.Mongo = &ir.IRMongoDBConfig{
					VersionKey:    m.Config.DB.Mongo.VersionKey,
					Collection:    m.Config.DB.Mongo.Collection,
					ToJSONGetters: m.Config.DB.Mongo.ToJSONGetters,
					Minimize:      m.Config.DB.Mongo.Minimize,
					AutoIndex:     m.Config.DB.Mongo.AutoIndex,
					AutoCreate:    m.Config.DB.Mongo.AutoCreate,
					StrictQuery:   m.Config.DB.Mongo.StrictQuery,
				}
				// mysql
				dbcfg.Mysql = &ir.IRMySQLDBConfig{
					Engine:  m.Config.DB.Mysql.Engine,
					Collate: m.Config.DB.Mysql.Collate,
				}
				cfg.DB = dbcfg
			}
			model.Config = cfg
		}

		// relations
		if m.Relations != nil {
			rel := &ir.IRRelations{}
			for _, mm := range m.Relations.ManyToMany {
				rel.ManyToMany = append(rel.ManyToMany, ir.IRManyToManyBlock{
					Name:     mm.Name,
					Ref:      mm.Ref,
					OnDelete: mm.OnDelete,
					OnUpdate: mm.OnUpdate,
				})
			}
			for _, hm := range m.Relations.HasMany {
				rel.HasMany = append(rel.HasMany, ir.IRHasManyBlock{
					Name:     hm.Name,
					Ref:      hm.Ref,
					OnDelete: hm.OnDelete,
					OnUpdate: hm.OnUpdate,
				})
			}
			for _, b := range m.Relations.BelongsTo {
				rel.BelongsTo = append(rel.BelongsTo, ir.IRBelongsToBlock{
					Name:     b.Name,
					Ref:      b.Ref,
					OnDelete: b.OnDelete,
					OnUpdate: b.OnUpdate,
				})
			}
			model.Relations = rel
		}

		ctx.IR.Models[model.Name] = model
	}

	return nil
}
