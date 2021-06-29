package jsmu

import "encoding/json"

var (
	// DefaultMarshaller is the Marshaller used when no Marshaller is provided to MU.
	DefaultMarshaller = &JSONMarshaller{}
)

// Marshaller is the interface for marshalling data.
type Marshaller interface {
	// Marshal marshals v into the expected encoding.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal unmarshals data into v.
	Unmarshal(data []byte, v interface{}) error
}

// JSONMarshaller implements Marshaller with JSON encoding.
type JSONMarshaller struct{}

// Marshal marshals v into the expected encoding.
func (me *JSONMarshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal unmarshals data into v.
func (me *JSONMarshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// MockMarshaller implements Marshaller but allows you to swap out the implementation with functions.
//
// Consider using this during unit testing.
type MockMarshaller struct {
	MarshalImplementation   func(v interface{}) ([]byte, error)
	UnmarshalImplementation func(data []byte, v interface{}) error
}

// Marshal marshals v into the expected encoding.
func (me *MockMarshaller) Marshal(v interface{}) ([]byte, error) {
	return me.MarshalImplementation(v)
}

// Unmarshal unmarshals data into v.
func (me *MockMarshaller) Unmarshal(data []byte, v interface{}) error {
	return me.UnmarshalImplementation(data, v)
}
