package jsmu_test

import (
	"fmt"

	"github.com/nofeaturesonlybugs/jsmu"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type ApiGetPerson struct {
	jsmu.TypeName `json:"-" jsmu:"get/person"`
	//
	Id     int     `json:"id"`
	Person *Person `json:"person"`
}

func ExampleMU_apiRequest() {
	//
	people := map[int]*Person{
		10: {
			Name: "Bob",
			Age:  40,
		},
		20: {
			Name: "Sally",
			Age:  32,
		},
	}
	//
	mu := &jsmu.MU{}
	mu.MustRegister(&ApiGetPerson{})
	//
	strings := []string{
		`{
			"type" : "get/person",
			"id" : "first",
			"message" : {
				"id" : 20
			}
		}`,
		`{
			"type" : "get/person",
			"id" : "second",
			"message" : {
				"id" : 10
			}
		}`,
	}
	var envelope jsmu.Enveloper
	var response []byte
	var err error
	for _, str := range strings {
		if envelope, err = mu.Unmarshal([]byte(str)); err != nil {
			fmt.Println(err)
			return
		}
		switch message := envelope.GetMessage().(type) {
		case *ApiGetPerson:
			if person, ok := people[message.Id]; ok {
				message.Person = person
			}
		}
		if response, err = mu.Marshal(envelope); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", string(response))
	}

	// Output: {"type":"get/person","id":"first","message":{"id":20,"person":{"name":"Sally","age":32}}}
	// {"type":"get/person","id":"second","message":{"id":10,"person":{"name":"Bob","age":40}}}
}
