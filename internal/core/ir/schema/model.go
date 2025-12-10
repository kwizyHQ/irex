package schema

// Model represents a parsed model from HCL
type Model struct {
	Name      string
	Fields    []Field
	Config    Config
	Relations []Relation
}

// Field represents a model field
type Field struct {
	Name        string
	Type        string
	Required    bool
	Optional    bool
	Unique      bool
	Default     interface{}
	Min         *float64
	Max         *float64
	MinLength   *int
	MaxLength   *int
	Match       string
	Trim        bool
	Visibility  string
	Description string
	Fields      []Field // nested object fields
}

// Config holds model-level configuration
type Config struct {
	Table       string
	Timestamps  bool
	IDStrategy  string
	Description string
	Strict      bool
	SoftDelete  bool
	Indexes     []Index
	DB          map[string]map[string]interface{}
}

// Index represents an index definition
type Index struct {
	Fields  []string
	Unique  bool
	Name    string
	Type    string
	Options map[string]interface{}
}

// Relation represents model relations
type Relation struct {
	Name         string
	Ref          string
	Type         string
	LocalField   string
	ForeignField string
	Through      string
	OnDelete     string
	OnUpdate     string
	Embedded     bool
}
