package ir

type DataAction string

const (
	DataCreate DataAction = "create"
	DataRead   DataAction = "read"
	DataUpdate DataAction = "update"
	DataDelete DataAction = "delete"
	DataList   DataAction = "list"
)

type OperationKind string

const (
	OperationKindData   OperationKind = "data"
	OperationKindCustom OperationKind = "custom"
)

type DataOperationMeta struct {
	Action DataAction `json:"action"`

	// cardinality
	Target string `json:"target"` // "single" | "many"

	// behavior
	Paginated  bool `json:"paginated,omitempty"`
	SoftDelete bool `json:"soft_delete,omitempty"`

	// semantics
	ReturnsEntity bool `json:"returns_entity,omitempty"`
	ReturnsList   bool `json:"returns_list,omitempty"`

	// optional hints
	OwnerField string `json:"owner_field,omitempty"`
}

type IROperation struct {
	Name        string `json:"name"`
	Service     string `json:"service,omitempty"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Action      string `json:"action,omitempty"`
	Description string `json:"description,omitempty"`

	Kind OperationKind `json:"kind"`

	Data *DataOperationMeta `json:"data,omitempty"`
}

type IROperations map[string]IROperation
