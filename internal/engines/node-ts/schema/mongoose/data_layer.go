package mongoose

import (
	"github.com/kwizyHQ/irex/internal/ir"
	"github.com/zclconf/go-cty/cty"
)

type MongoModel struct {
	Name        string
	Fields      []MongoField
	Indexes     []MongoIndex
	Config      MongoModelConfig
	Relations   []MongoRelation
	Description string
}

type MongoField struct {
	Name         string
	MongooseType string
	IsArray      bool
	Required     bool
	Unique       bool
	Trim         bool
	Min          *int
	Max          *int
	MinLength    *int
	MaxLength    *int
	Match        string
	Default      cty.Value
	Index        bool
	Collation    *MongoCollation
	SubFields    []MongoField
	Description  string
	Visibility   string
}

type MongoCollation struct {
	Locale          string
	CaseLevel       bool
	CaseFirst       string
	Strength        int
	NumericOrdering bool
	Alternate       string
	MaxVariable     string
	Backwards       bool
}

type MongoIndex struct {
	Name   string
	Fields []string
	Unique bool
}

type MongoModelConfig struct {
	Collection    string
	Timestamps    bool
	Strict        bool
	VersionKey    bool
	AutoIndex     bool
	AutoCreate    bool
	Minimize      bool
	StrictQuery   bool
	ToJSONGetters bool
}

type MongoRelation struct {
	Name     string
	Ref      string
	Type     string // hasMany | belongsTo | manyToMany
	OnDelete string
	OnUpdate string
}

func BuildMongoModel(m ir.IRModel) MongoModel {
	model := MongoModel{
		Name: m.Name,
		// Description: m.Config.Description,
	}

	// Fields
	for _, f := range m.Fields {
		model.Fields = append(model.Fields, buildMongoField(f))
	}

	// Indexes
	if m.Config != nil {
		for _, idx := range m.Config.Indexes {
			model.Indexes = append(model.Indexes, MongoIndex{
				Name:   idx.Name,
				Fields: idx.Fields,
				Unique: idx.Unique,
			})
		}
	}

	// Config
	if m.Config != nil && m.Config.DB != nil && m.Config.DB.Mongo != nil {
		db := m.Config.DB.Mongo
		model.Config = MongoModelConfig{
			Collection:    db.Collection,
			Timestamps:    m.Config.Timestamps,
			Strict:        m.Config.Strict,
			VersionKey:    db.VersionKey,
			AutoIndex:     db.AutoIndex,
			AutoCreate:    db.AutoCreate,
			Minimize:      db.Minimize,
			StrictQuery:   db.StrictQuery,
			ToJSONGetters: db.ToJSONGetters,
		}
	}

	// Relations
	if m.Relations != nil {
		for _, r := range m.Relations.HasMany {
			model.Relations = append(model.Relations, MongoRelation{
				Name: r.Name,
				Ref:  r.Ref,
				Type: "hasMany",
			})
		}
		for _, r := range m.Relations.BelongsTo {
			model.Relations = append(model.Relations, MongoRelation{
				Name: r.Name,
				Ref:  r.Ref,
				Type: "belongsTo",
			})
		}
		for _, r := range m.Relations.ManyToMany {
			model.Relations = append(model.Relations, MongoRelation{
				Name: r.Name,
				Ref:  r.Ref,
				Type: "manyToMany",
			})
		}
	}

	return model
}

func buildMongoField(f ir.IRModelField) MongoField {
	field := MongoField{
		Name:        f.Name,
		Required:    f.Required,
		Unique:      f.Unique,
		Trim:        f.Trim,
		Min:         f.Min,
		Max:         f.Max,
		MinLength:   f.MinLength,
		MaxLength:   f.MaxLength,
		Match:       f.Match,
		Description: f.Description,
		Visibility:  f.Visibility,
	}

	field.MongooseType, field.IsArray = mapMongooseType(f.Type)

	// Embedded fields
	for _, sf := range f.Fields {
		field.SubFields = append(field.SubFields, buildMongoField(sf))
	}

	// Mongo DB config
	if f.DB != nil && f.DB.Mongo != nil {
		field.Index = f.DB.Mongo.Index
		field.Unique = field.Unique || f.DB.Mongo.Unique
		if f.DB.Mongo.Collation != nil {
			field.Collation = &MongoCollation{
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

	field.Default = f.Default

	// Default value (stringified)
	// if f.Default.IsKnown() && !f.Default.IsNull() {
	// 	field.Default = fmt.Sprintf("%v", f.Default.GoString())
	// }

	return field
}

func mapMongooseType(t string) (string, bool) {
	switch t {
	case "string":
		return "String", false
	case "int", "float", "number":
		return "Number", false
	case "bool":
		return "Boolean", false
	case "date":
		return "Date", false
	case "objectId":
		return "mongoose.Schema.Types.ObjectId", false
	case "string[]":
		return "String", true
	default:
		return "mongoose.Schema.Types.Mixed", false
	}
}
