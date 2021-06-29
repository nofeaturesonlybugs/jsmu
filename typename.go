package jsmu

import "reflect"

// TypeName represents a jsmu type name.  Embed the TypeName type into a struct and set
// the approprate struct tag to configure the table name.
type TypeName string

var (
	typeTypeName = reflect.TypeOf(TypeName(""))
)
