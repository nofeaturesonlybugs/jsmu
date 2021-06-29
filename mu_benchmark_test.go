package jsmu_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/nofeaturesonlybugs/jsmu"
	"github.com/nofeaturesonlybugs/jsmu/data"
)

func BenchmarkMU(b *testing.B) {
	people := strings.Split(strings.Trim(data.PeopleStream, "\r\n "), "\n")
	animals := strings.Split(strings.Trim(data.AnimalStream, "\r\n "), "\n")
	all := make([]string, len(people)+len(animals))
	copy(all, people)
	copy(all[len(people):], animals)
	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})
	//
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
		// EnveloperFn: nil,  // jsmu.DefaultEnveloperFunc is the default.
		// StructTag  : "",   // "jsmu" is the default.
	}
	mu.MustRegister(&Person{})
	mu.MustRegister(&Animal{})
	//
	limits := []int{
		5,
		100,
		250,
		500,
		1000,
	}
	for _, limit := range limits {
		b.Run(fmt.Sprintf("2-pass unmarshal limit %v", limit), func(b *testing.B) {
			type E struct {
				Type    string          `json:"type"`
				Id      string          `json:"id"`
				Message json.RawMessage `json:"message"`
				V       interface{}     `json:"-"`
			}
			var err error
			for k := 0; k < b.N; k++ {
				for n, max := 0, 2*limit; n < max; n++ {
					env := E{}
					if err = json.Unmarshal([]byte(all[n]), &env); err != nil {
						b.Fatalf("during first json.Unmarshal: %v", err.Error())
					}
					switch env.Type {
					case "animal":
						env.V = &Animal{}
					case "person":
						env.V = &Person{}
					}
					if err = json.Unmarshal(env.Message, env.V); err != nil {
						b.Fatalf("during second json.Unmarshal: %v", err.Error())
					}
				}
			}
		})
		b.Run(fmt.Sprintf("jsmu unmarshal limit %v", limit), func(b *testing.B) {
			var err error
			for k := 0; k < b.N; k++ {
				for n, max := 0, 2*limit; n < max; n++ {
					if _, err = mu.Unmarshal([]byte(all[n])); err != nil {
						b.Fatalf("during MU.Unmarshal: %v", err.Error())
					}
				}
			}
		})
	}
}
