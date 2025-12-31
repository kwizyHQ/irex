package ir

import "github.com/zclconf/go-cty/cty"

// IR model types mirror the resolved model/schema definitions from
// internal/core/symbols and provide a generator-friendly representation.

type IRMongoFieldCollation struct {
	Locale          string `json:"locale,omitempty"`
	CaseLevel       bool   `json:"case_level,omitempty"`
	CaseFirst       string `json:"case_first,omitempty"`
	Strength        int    `json:"strength,omitempty"`
	NumericOrdering bool   `json:"numeric_ordering,omitempty"`
	Alternate       string `json:"alternate,omitempty"`
	MaxVariable     string `json:"max_variable,omitempty"`
	Backwards       bool   `json:"backwards,omitempty"`
}

type IRMongoFieldConfig struct {
	Index     bool                   `json:"index,omitempty"`
	Unique    bool                   `json:"unique,omitempty"`
	Collation *IRMongoFieldCollation `json:"collation,omitempty"`
}

type IRMySQLFieldConfig struct {
	Index   bool   `json:"index,omitempty"`
	Unique  bool   `json:"unique,omitempty"`
	Collate string `json:"collate,omitempty"`
}

type IRModelFieldDBConfig struct {
	Mongo *IRMongoFieldConfig `json:"mongo,omitempty"`
	Mysql *IRMySQLFieldConfig `json:"mysql,omitempty"`
}

type IRModelField struct {
	Name        string                `json:"name"`
	Type        string                `json:"type,omitempty"`
	Required    bool                  `json:"required,omitempty"`
	Unique      bool                  `json:"unique,omitempty"`
	Trim        bool                  `json:"trim,omitempty"`
	MinLength   *int                  `json:"min_length,omitempty"`
	MaxLength   *int                  `json:"max_length,omitempty"`
	Min         *int                  `json:"min,omitempty"`
	Max         *int                  `json:"max,omitempty"`
	Default     cty.Value             `json:"default,omitempty"`
	Match       string                `json:"match,omitempty"`
	Message     string                `json:"message,omitempty"`
	Visibility  string                `json:"visibility,omitempty"`
	Fields      []IRModelField        `json:"fields,omitempty"`
	DB          *IRModelFieldDBConfig `json:"db,omitempty"`
	Description string                `json:"description,omitempty"`
}

type IRModelIndex struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields,omitempty"`
	Unique bool     `json:"unique,omitempty"`
}

type IRMongoDBConfig struct {
	VersionKey    bool   `json:"version_key,omitempty"`
	Collection    string `json:"collection,omitempty"`
	ToJSONGetters bool   `json:"to_json_getters,omitempty"`
	Minimize      bool   `json:"minimize,omitempty"`
	AutoIndex     bool   `json:"auto_index,omitempty"`
	AutoCreate    bool   `json:"auto_create,omitempty"`
	StrictQuery   bool   `json:"strict_query,omitempty"`
}

type IRMySQLDBConfig struct {
	Engine  string `json:"engine,omitempty"`
	Collate string `json:"collate,omitempty"`
}

type IRModelConfigDB struct {
	Mongo *IRMongoDBConfig `json:"mongo,omitempty"`
	Mysql *IRMySQLDBConfig `json:"mysql,omitempty"`
}

type IRModelConfig struct {
	Timestamps  bool             `json:"timestamps,omitempty"`
	Table       string           `json:"table,omitempty"`
	Strict      bool             `json:"strict,omitempty"`
	Indexes     []IRModelIndex   `json:"indexes,omitempty"`
	IDStrategy  string           `json:"id_strategy,omitempty"`
	Description string           `json:"description,omitempty"`
	DB          *IRModelConfigDB `json:"db,omitempty"`
}

type IRManyToManyBlock struct {
	Name     string `json:"name"`
	Ref      string `json:"ref"`
	OnDelete string `json:"on_delete,omitempty"`
	OnUpdate string `json:"on_update,omitempty"`
}

type IRHasManyBlock struct {
	Name     string `json:"name"`
	Ref      string `json:"ref"`
	OnDelete string `json:"on_delete,omitempty"`
	OnUpdate string `json:"on_update,omitempty"`
}

type IRBelongsToBlock struct {
	Name     string `json:"name"`
	Ref      string `json:"ref"`
	OnDelete string `json:"on_delete,omitempty"`
	OnUpdate string `json:"on_update,omitempty"`
}

type IRRelations struct {
	ManyToMany []IRManyToManyBlock `json:"many_to_many,omitempty"`
	HasMany    []IRHasManyBlock    `json:"has_many,omitempty"`
	BelongsTo  []IRBelongsToBlock  `json:"belongs_to,omitempty"`
}

type IRModel struct {
	Name      string         `json:"name"`
	Fields    []IRModelField `json:"fields,omitempty"`
	Config    *IRModelConfig `json:"config,omitempty"`
	Relations *IRRelations   `json:"relations,omitempty"`
}

type IRModels map[string]IRModel
