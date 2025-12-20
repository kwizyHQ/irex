package symbols

import (
	"github.com/zclconf/go-cty/cty"
)

type MongoDBFieldCollation struct {
	Locale          string `hcl:"locale,optional"`
	CaseLevel       bool   `hcl:"caseLevel,optional"`
	CaseFirst       string `hcl:"caseFirst,optional"`
	Strength        int    `hcl:"strength,optional"`
	NumericOrdering bool   `hcl:"numericOrdering,optional"`
	Alternate       string `hcl:"alternate,optional"`
	MaxVariable     string `hcl:"maxVariable,optional"`
	Backwards       bool   `hcl:"backwards,optional"`
}

type MongoDBFieldConfig struct {
	Index     bool                   `hcl:"index,optional"`
	Unique    bool                   `hcl:"unique,optional"`
	Collation *MongoDBFieldCollation `hcl:"collation,block"`
}

type MySqlDBFieldConfig struct {
	Index   bool   `hcl:"index,optional"`
	Unique  bool   `hcl:"unique,optional"`
	Collate string `hcl:"collate,optional"`
}

// Model definitions generated from models.hcl (HCL tags)
type ModelFieldDBConfig struct {
	Mongo *MongoDBFieldConfig `hcl:"mongo,block"`
	Mysql *MySqlDBFieldConfig `hcl:"mysql,block"`
}

type ModelField struct {
	Name        string              `hcl:"name,label"`
	Type        string              `hcl:"type,optional"`
	Required    bool                `hcl:"required,optional"`
	Unique      bool                `hcl:"unique,optional"`
	Trim        bool                `hcl:"trim,optional"`
	MinLength   *int                `hcl:"minlength,optional"`
	MaxLength   *int                `hcl:"maxlength,optional"`
	Min         *int                `hcl:"min,optional"`
	Max         *int                `hcl:"max,optional"`
	Default     cty.Value           `hcl:"default,optional"`
	Match       string              `hcl:"match,optional"`
	Message     string              `hcl:"message,optional"`
	Visibility  string              `hcl:"visibility,optional"`
	Fields      []ModelField        `hcl:"field,block"` // for nested fields
	DB          *ModelFieldDBConfig `hcl:"db,block"`
	Description string              `hcl:"description,optional"`
}

type ModelIndex struct {
	Name   string   `hcl:"name,label"`
	Fields []string `hcl:"fields"`
	Unique bool     `hcl:"unique,optional"`
}

type MongoDBConfig struct {
	VersionKey    bool   `hcl:"versionKey,optional"`
	Collection    string `hcl:"collection,optional"`
	ToJSONGetters bool   `hcl:"toJSON_getters,optional"`
	Minimize      bool   `hcl:"minimize,optional"`
	AutoIndex     bool   `hcl:"autoIndex,optional"`
	AutoCreate    bool   `hcl:"autoCreate,optional"`
	StrictQuery   bool   `hcl:"strictQuery,optional"`
}

type MySqlDBConfig struct {
	Engine  string `hcl:"engine,optional"`
	Collate string `hcl:"collate,optional"`
}

type ModelConfigDB struct {
	Mongo MongoDBConfig `hcl:"mongo,block"`
	Mysql MySqlDBConfig `hcl:"mysql,block"`
}

type ModelConfig struct {
	Timestamps  bool           `hcl:"timestamps,optional"`
	Table       string         `hcl:"table,optional"`
	Strict      bool           `hcl:"strict,optional"`
	Indexes     []ModelIndex   `hcl:"index,block"`
	IDStrategy  string         `hcl:"idStrategy,optional"`
	Description string         `hcl:"description,optional"`
	DB          *ModelConfigDB `hcl:"db,block"`
}

type ManyToManyBlock struct {
	Name     string `hcl:"name,label"`
	Ref      string `hcl:"ref"`
	OnDelete string `hcl:"onDelete,optional"`
	OnUpdate string `hcl:"onUpdate,optional"`
}

type HasManyBlock struct {
	Name     string `hcl:"name,label"`
	Ref      string `hcl:"ref"`
	OnDelete string `hcl:"onDelete,optional"`
	OnUpdate string `hcl:"onUpdate,optional"`
}

type BelongsToBlock struct {
	Name     string `hcl:"name,label"`
	Ref      string `hcl:"ref"`
	OnDelete string `hcl:"onDelete,optional"`
	OnUpdate string `hcl:"onUpdate,optional"`
}

type Relations struct {
	ManyToMany []ManyToManyBlock `hcl:"manyToMany,block"`
	HasMany    []HasManyBlock    `hcl:"hasMany,block"`
	BelongsTo  []BelongsToBlock  `hcl:"belongsTo,block"`
}

type Model struct {
	Name      string       `hcl:"name,label"`
	Fields    []ModelField `hcl:"field,block"`
	Config    *ModelConfig `hcl:"config,block"`
	Relations *Relations   `hcl:"relations,block"`
}

type ModelsBlock struct {
	Models []Model `hcl:"model,block"`
}

type ModelsSpec struct {
	ModelsBlock *ModelsBlock `hcl:"models,block"`
}
