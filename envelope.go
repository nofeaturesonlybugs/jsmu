package jsmu

import "encoding/json"

var (
	// DefaultEnveloperFunc creates and returns an *Envelope from this package.
	DefaultEnveloperFunc = func() Enveloper {
		return &Envelope{}
	}
)

// EnveloperFunc is a function that creates a new Enveloper.
type EnveloperFunc func() Enveloper

// Enveloper is the interface for message Envelopes.
type Enveloper interface {
	// GetMessage returns the message in the envelope.
	GetMessage() interface{}
	// SetMessage sets the message in the envelope.
	SetMessage(interface{})
	// GetRawMessage returns the raw JSON message in the envelope.
	GetRawMessage() json.RawMessage
	// SetRawMessage sets the raw JSON message in the envelope.
	SetRawMessage(json.RawMessage)
	// GetTypeName returns the type name string.
	GetTypeName() string
	// SetTypeName sets the type name string.
	SetTypeName(string)
}

// Envelope is a default Enveloper for basic needs.
//
// JSON must conform to this format:
//	{
//		"type"       : string,  // A string uniquely identifying "message".
//		"id"         : string,  // An optional unique message ID; provided as a convenience for request-reply message flow.
//		"message"    : {},      // Any JSON data your application can receive.
//	}
type Envelope struct {
	Type string          `json:"type"`
	Id   string          `json:"id"`
	Raw  json.RawMessage `json:"message"`
	//
	// This is the concrete type; even though it's a private field and invisible to the json
	// package we still explicitly exclude it.
	message interface{} `json:"-"`
}

// GetMessage returns the message in the envelope.
func (me *Envelope) GetMessage() interface{} {
	return me.message
}

// SetMessage sets the message in the envelope.
func (me *Envelope) SetMessage(message interface{}) {
	me.message = message
}

// GetRawMessage returns the raw JSON message in the envelope.
func (me *Envelope) GetRawMessage() json.RawMessage {
	return me.Raw
}

// SetRawMessage sets the raw JSON message in the envelope.
func (me *Envelope) SetRawMessage(raw json.RawMessage) {
	me.Raw = raw
}

// GetTypeName returns the type name string.
func (me *Envelope) GetTypeName() string {
	return me.Type
}

// SetTypeName sets the type name string.
func (me *Envelope) SetTypeName(typeName string) {
	me.Type = typeName
}
