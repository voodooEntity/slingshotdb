package types

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// - - - - - - - - STORAGE STRUCTS - - - - - - - - - -
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// StorageEntity struct
type StorageEntity struct {
	ID         int
	Type       int
	Context    string
	Value      string
	Properties map[string]string
	Version    int
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation struct
type StorageRelation struct {
	SourceType int
	SourceID   int
	TargetType int
	TargetID   int
	Context    string
	Properties map[string]string
	Version    int
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// - - - - - - - - MAPPER STRUCTS  - - - - - - - - - -
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// entity transport struct for json<->entity
type MapperEntity struct {
	ID         int
	Type       string
	Context    string
	Value      string
	Properties map[string]string
	Children   []MapperEntity
	Version    int
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// relation transport struct
type MapperRelation struct {
	SourceType string
	SourceID   int
	TargetType string
	TargetID   int
	Context    string
	Properties map[string]string
	Version    int
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// mapper transport  struct
type MapperTransport struct {
	Entities  []MapperEntity
	Relations []MapperRelation
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// - - - - - - PERSISTANCE STRUCTS - - - - - - - - - -
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// persistance struct
type PersistancePayload struct {
	Type        string
	Method      string
	Entity      StorageEntity
	Relation    StorageRelation
	EntityTypes map[int]string
}

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// - - - - - - - -  CONFIG STRUCTS - - - - - - - - - -
// - - - - - - - - - - - - - - - - - - - - - - - - - -

// - - - - - - - - - - - - - - - - - - - - - - - - - -
// config struct
type Config struct {
	Host        string
	Port        int
	Persistance bool
}
