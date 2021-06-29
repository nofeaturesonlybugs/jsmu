package jsmu_test

// This example demonstrates using a custom envelope with a JSON structure not expected by jsmu.Envelope.

import (
	"encoding/json"
	"fmt"

	"github.com/nofeaturesonlybugs/jsmu"
)

// CustomEnvelope is for JSON that has a custom structure.
type CustomEnvelope struct {
	TypeInfo string          `json:"typeInfo"`
	Payload  json.RawMessage `json:"payload"`
	CustomA  string          `json:"customA"`
	CustomB  string          `json:"customB"`
	data     interface{}     `json:"-"`
}

// GetMessage returns the message in the envelope.
func (me *CustomEnvelope) GetMessage() interface{} {
	return me.data
}

// SetMessage sets the message in the envelope.
func (me *CustomEnvelope) SetMessage(message interface{}) {
	me.data = message
}

// GetRawMessage returns the raw JSON message in the envelope.
func (me *CustomEnvelope) GetRawMessage() json.RawMessage {
	return me.Payload
}

// SetRawMessage sets the raw JSON message in the envelope.
func (me *CustomEnvelope) SetRawMessage(raw json.RawMessage) {
	me.Payload = raw
}

// GetTypeName returns the type name string.
func (me *CustomEnvelope) GetTypeName() string {
	return me.TypeInfo
}

// SetTypeName sets the type name string.
func (me *CustomEnvelope) SetTypeName(typeInfo string) {
	me.TypeInfo = typeInfo
}

func ExampleMU_customEnvelope() {
	type Person struct {
		jsmu.TypeName `js:"-" jsmu:"person"`
		//
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Animal struct {
		jsmu.TypeName `js:"-" jsmu:"animal"`
		//
		Name string `json:"name"`
		Says string `json:"says"`
	}
	//
	mu := &jsmu.MU{
		EnveloperFn: func() jsmu.Enveloper {
			return &CustomEnvelope{}
		},
		// StructTag  : "",   // "jsmu" is the default.
	}
	mu.MustRegister(&Person{})
	mu.MustRegister(&Animal{})
	//
	strings := []string{
		`{
			"typeInfo" : "person",
			"payload" : {
				"name" : "Bob",
				"age" : 40
			},
			"customA" : "Hello!",
			"customB" : "Goodbye!"
		}`,
		`{
			"typeInfo" : "animal",
			"payload" : {
				"name" : "cat",
				"says" : "meow"
			},
			"customA" : "Foo",
			"customB" : "Bar"
		}`,
		`{
			"typeInfo" : "person",
			"payload" : {
				"name" : "Sally",
				"age" : 30
			},
			"customA" : "qwerty",
			"customB" : "dvorak"
		}`,
		`{
			"typeInfo" : "animal",
			"payload" : {
				"name" : "cow",
				"says" : "moo"
			},
			"customA" : "Batman",
			"customB" : "Robin"
		}`,
	}
	var envelope jsmu.Enveloper
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		custom := envelope.(*CustomEnvelope)
		switch message := envelope.GetMessage().(type) {
		case *Animal:
			fmt.Printf("A %v says, \"%v;\" [%v,%v].\n", message.Name, message.Says, custom.CustomA, custom.CustomB)
		case *Person:
			fmt.Printf("%v is %v year(s) old; [%v,%v].\n", message.Name, message.Age, custom.CustomA, custom.CustomB)
		}
	}

	// Output: Bob is 40 year(s) old; [Hello!,Goodbye!].
	// A cat says, "meow;" [Foo,Bar].
	// Sally is 30 year(s) old; [qwerty,dvorak].
	// A cow says, "moo;" [Batman,Robin].
}
