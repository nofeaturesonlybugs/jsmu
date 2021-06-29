package jsmu_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/jsmu"
)

func ExampleMU() {
	// This example demonstrates:
	//	+ Register() when the type information is embedded with jsmu.TypeName.
	//	+ Unmarshal()

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
		// EnveloperFn: nil,  // jsmu.DefaultEnveloperFunc is the defualt.
		// StructTag  : "",   // "jsmu" is the default.
	}
	mu.MustRegister(&Person{})
	mu.MustRegister(&Animal{})
	//
	strings := []string{
		`{
			"type" : "person",
			"message" : {
				"name" : "Bob",
				"age" : 40
			}
		}`,
		`{
			"type" : "animal",
			"message" : {
				"name" : "cat",
				"says" : "meow"
			}
		}`,
		`{
			"type" : "person",
			"message" : {
				"name" : "Sally",
				"age" : 30
			}
		}`,
		`{
			"type" : "animal",
			"message" : {
				"name" : "cow",
				"says" : "moo"
			}
		}`,
	}
	var envelope jsmu.Enveloper
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		switch message := envelope.GetMessage().(type) {
		case *Animal:
			fmt.Printf("A %v says, \"%v.\"\n", message.Name, message.Says)
		case *Person:
			fmt.Printf("%v is %v year(s) old.\n", message.Name, message.Age)
		}
	}

	// Output: Bob is 40 year(s) old.
	// A cat says, "meow."
	// Sally is 30 year(s) old.
	// A cow says, "moo."
}

func ExampleMU_withoutEmbed() {
	// This example does not embed jsmu.TypeName into the registered structs and
	// instead passes that information in the call to Register().
	//
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Animal struct {
		Name string `json:"name"`
		Says string `json:"says"`
	}
	//
	mu := &jsmu.MU{}
	mu.MustRegister(&Person{}, jsmu.TypeName("person"))
	mu.MustRegister(&Animal{}, jsmu.TypeName("animal"))
	//
	strings := []string{
		`{
			"type" : "person",
			"message" : {
				"name" : "Bob",
				"age" : 40
			}
		}`,
		`{
			"type" : "animal",
			"message" : {
				"name" : "cat",
				"says" : "meow"
			}
		}`,
		`{
			"type" : "person",
			"message" : {
				"name" : "Sally",
				"age" : 30
			}
		}`,
		`{
			"type" : "animal",
			"message" : {
				"name" : "cow",
				"says" : "moo"
			}
		}`,
	}
	var envelope jsmu.Enveloper
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		switch message := envelope.GetMessage().(type) {
		case *Animal:
			fmt.Printf("A %v says, \"%v.\"\n", message.Name, message.Says)
		case *Person:
			fmt.Printf("%v is %v year(s) old.\n", message.Name, message.Age)
		}
	}

	// Output: Bob is 40 year(s) old.
	// A cat says, "meow."
	// Sally is 30 year(s) old.
	// A cow says, "moo."
}

func ExampleMU_messageFactory() {
	// This example shows how to provide a constructor or factory function for instantiating
	// your types if necessary.

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	// We create a constructor function compatible with jsmu.MessageFunc.
	NewPerson := func() (interface{}, error) {
		fmt.Println("Created a person!")
		return &Person{}, nil
	}
	//
	mu := &jsmu.MU{
		// EnveloperFn: nil,  // jsmu.DefaultEnveloperFunc is the defualt.
		// StructTag  : "",   // "jsmu" is the default.
	}
	// Pass the NewPerson constructor during registration.
	mu.MustRegister(&Person{}, jsmu.MessageFunc(NewPerson), jsmu.TypeName("person"))
	//
	strings := []string{
		`{
			"type" : "person",
			"message" : {
				"name" : "Bob",
				"age" : 40
			}
		}`,
		`{
			"type" : "person",
			"message" : {
				"name" : "Sally",
				"age" : 30
			}
		}`,
	}
	var envelope jsmu.Enveloper
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		switch message := envelope.GetMessage().(type) {
		case *Person:
			fmt.Printf("%v is %v year(s) old.\n", message.Name, message.Age)
		}
	}

	// Output: Created a person!
	// Bob is 40 year(s) old.
	// Created a person!
	// Sally is 30 year(s) old.
}

func ExampleMU_arrays() {
	// This example shows how to register slices.

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type Animal struct {
		Name string `json:"name"`
		Says string `json:"says"`
	}
	//
	mu := &jsmu.MU{}
	mu.MustRegister([]*Person{}, jsmu.TypeName("people"))
	mu.MustRegister([]*Animal{}, jsmu.TypeName("animals"))
	//
	strings := []string{
		`{
			"type" : "people",
			"message" : [{
				"name" : "Bob",
				"age" : 40
			}, {
				"name" : "Sally",
				"age" : 30
			}]
		}`,
		`{
			"type" : "animals",
			"message" : [{
				"name" : "cat",
				"says" : "meow"
			}, {
				"name" : "cow",
				"says" : "moo"
			}]
		}`,
	}
	var envelope jsmu.Enveloper
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		switch message := envelope.GetMessage().(type) {
		case []*Animal:
			for _, animal := range message {
				fmt.Printf("A %v says, \"%v.\"\n", animal.Name, animal.Says)
			}
		case []*Person:
			for _, person := range message {
				fmt.Printf("%v is %v year(s) old.\n", person.Name, person.Age)
			}
		}
	}

	var buf []byte
	people := []*Person{
		{Name: "Larry", Age: 42},
		{Name: "Curly", Age: 48},
		{Name: "Moe", Age: 46},
	}
	if buf, err = mu.Marshal(people); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", string(buf))

	// Output: Bob is 40 year(s) old.
	// Sally is 30 year(s) old.
	// A cat says, "meow."
	// A cow says, "moo."
	// {"type":"people","id":"","message":[{"name":"Larry","age":42},{"name":"Curly","age":48},{"name":"Moe","age":46}]}
}
