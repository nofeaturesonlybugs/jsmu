package jsmu

import "fmt"

var (
	// ErrAlreadyRegistered is returned during Register() if the type string is already registered.
	ErrAlreadyRegistered = fmt.Errorf("already registered")
	// ErrMissingTypeName is returned during Register() if the type string can not be located.
	ErrMissingTypeName = fmt.Errorf("missing type name")
	// ErrNothingRegistered is returned during MU.Marshal() and MU.Unmarshal() when no types have been registered.
	ErrNothingRegistered = fmt.Errorf("nothing registered")
	// ErrUnregisteredType is returned during MU.Marshal() and MU.Unmarshal() when the type has not been registered.
	ErrUnregisteredType = fmt.Errorf("unregistered type")
)
