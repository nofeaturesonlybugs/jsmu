// Package jsmu abstracts the typical 2-pass JSON encoding into what looks like a 1-pass JSON encoding step when
// marshalling or unmarshalling data.
//
// Pronounce as j◦s◦mew.
//
// TypeName, ConcreteType, and 2-Pass Encoding
//
// Writing a general purpose JSON unmarshaller usually requires the following steps:
//	1. json.Unmarshal() into a type containing at least two fields:
//		a. int|string field identifying the desired Go struct destination type; call this the TypeName field.
//		b. json.RawMessage field containing the data or payload that is not yet unmarshalled.
//	2. Type switch or branch off the TypeName field to create a Go struct destination; call this the ConcreteType
//	3. json.Unmarshal() the json.RawMessage into ConcreteType.
//
// Writing a general purpose JSON marshaller works in reverse:
//	1. json.Marshal() the outgoing ConcreteType to json.RawMessage.
//	2. Type switch or branch off the outgoing ConcreteType to create an int|string TypeName identification.
//	3. Pack the outputs from steps 1 and 2 (along with any other desired information) into some type of envelope and json.Marshal() the envelope.
//
// jsmu simply abstracts the above logic into a single marshal or unmarshal call as necessary and alleviates you from having
// to maintain large type switches, constructor tables, or whatever other mechanism you perform the above logic with.
//
// Register
//
// jsmu requires that ConcreteType(s) be registered along with their TypeName and supports two conventions
// for supplying the TypeName.
//
// You can embed a TypeName field directly into your struct:
//	type MyType struct {
//		jsmu.TypeName `jsmu:"my-type"`
//		//
//		String string `json:"string"`
//		Number int `json:"number"`
//	}
//	var mu jsmu.MU
//	mu.Register(&MyType{})
//
// Or you can pass a TypeName value during the call to Register():
//	type MyType struct {
//		String string `json:"string"`
//		Number int `json:"number"`
//	}
//	//
//	var mu jsmu.MU
//	mu.Register(&MyType{}, jsmu.TypeName("my-type"))
//
// Messages and Envelopes
//
// Applications want to send messages.  In order for an endpoint to unpack an arbitrary message into a concrete type
// we need to send type information along with the message (this was referred to as TypeName above).
//
// In other words the most basic message format becomes:
//	{
//		"type"     : int|string     // Uniquely describes the type for the "message" field.
//		"message"  : any            // The message (broadcast|request|reply) the application wanted to send.
//	}
//
// Conceptually we can think of this type wrapped around the "message" as the envelope.
//
// jsmu allows you to use any custom type as an envelope by exposing and expecting the Enveloper interface:
//	type CustomEnvelope struct {
//		// Must implement jsmu.Enveloper interface
//	}
//	my := &jsmu.MU{
//		EnveloperFn : func() Enveloper{ return &CustomEnvelope{} },
//	}
//
// If you leave the EnveloperFn field as nil then DefaultEnveloperFunc is used and your JSON needs to conform to the structure
// defined by Envelope.
package jsmu
