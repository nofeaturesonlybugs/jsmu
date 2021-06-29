package jsmu_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/jsmu"
)

func TestMU_MustRegister_DoubleRegisterPanics(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct {
		jsmu.TypeName `jsmu:"a"`
		A             string
	}
	didPanic := false
	str := ""
	mu := jsmu.MU{}
	mu.MustRegister(&A{})
	//
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
				str = fmt.Sprintf("%v", r)
			}
		}()
		mu.MustRegister(&A{})
	}()
	chk.True(didPanic)
	chk.Contains(str, "already registered")
}

func TestMU_Register_NoTypeNameIsError(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct {
		A string
	}
	mu := jsmu.MU{}
	err := mu.Register(&A{})
	chk.Error(err)
	chk.Contains(err.Error(), "missing type name")
}

func TestMU_Marshal_NothingRegisteredIsError(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct{}
	mu := jsmu.MU{}
	//
	_, err := mu.Marshal(&A{})
	chk.Error(err)
	chk.Contains(err.Error(), "nothing registered")
}

func TestMU_Unmarshal_NothingRegisteredIsError(t *testing.T) {
	chk := assert.New(t)
	//
	js := `{
		"type" : "person",
		"message" : {
			"name" : "Bob",
			"age" : 40
		}
	}`
	mu := jsmu.MU{}
	//
	_, err := mu.Unmarshal([]byte(js))
	chk.Error(err)
	chk.Contains(err.Error(), "nothing registered")
}

func TestMU_Marshal_UnregisteredIsError(t *testing.T) {
	chk := assert.New(t)
	//
	type A struct{}
	type B struct{}
	mu := jsmu.MU{}
	// Need to register at least one thing for the other to be reported as unregistered.
	mu.Register(&B{}, jsmu.TypeName("B"))
	//
	_, err := mu.Marshal(&A{})
	chk.Error(err)
	chk.Contains(err.Error(), "unregistered type")
}

func TestMU_Unmarshal_UnregisteredIsError(t *testing.T) {
	chk := assert.New(t)
	//
	type B struct{}
	//
	js := `{
		"type" : "person",
		"message" : {
			"name" : "Bob",
			"age" : 40
		}
	}`
	mu := jsmu.MU{}
	// Need to register at least one thing for the other to be reported as unregistered.
	mu.Register(&B{}, jsmu.TypeName("B"))
	//
	_, err := mu.Unmarshal([]byte(js))
	chk.Error(err)
	chk.Contains(err.Error(), "unregistered type")
}

func TestMU_Unmarshal_CustomFactoryReturnsError(t *testing.T) {
	chk := assert.New(t)
	//
	type Person struct {
		jsmu.TypeName `jsmu:"person"`
		Name          string `json:"name"`
		Age           int    `json:"age"`
	}
	fn := func() (interface{}, error) {
		return nil, fmt.Errorf("custom factory error")
	}
	js := `{
		"type" : "person",
		"message" : {
			"name" : "Bob",
			"age" : 40
		}
	}`
	mu := jsmu.MU{}
	mu.MustRegister(&Person{}, jsmu.MessageFunc(fn))
	_, err := mu.Unmarshal([]byte(js))
	chk.Error(err)
}

func TestMU_Unmarshal_MarshallerReturnsErrors(t *testing.T) {
	chk := assert.New(t)
	//
	type Person struct {
		jsmu.TypeName `jsmu:"person"`
		Name          string `json:"name"`
		Age           int    `json:"age"`
	}
	js := `{
		"type" : "person",
		"message" : {
			"name" : "Bob",
			"age" : 40
		}
	}`
	//
	mock := &jsmu.MockMarshaller{}
	//
	mu := jsmu.MU{
		Marshaller: mock,
	}
	mu.MustRegister(&Person{})
	//
	{
		// First call to Unmarshal errors
		mock.UnmarshalImplementation = func(data []byte, v interface{}) error {
			return errors.Errorf("first unmarshal error")
		}
		_, err := mu.Unmarshal([]byte(js))
		chk.Error(err)
		chk.Contains(err.Error(), "first unmarshal error")
	}
	//
	{
		// Second call to Unmarshal errors
		calls := -1
		mock.UnmarshalImplementation = func(data []byte, v interface{}) error {
			calls++
			if calls == 0 {
				return jsmu.DefaultMarshaller.Unmarshal(data, v)
			}
			return errors.Errorf("second unmarshal error")
		}
		_, err := mu.Unmarshal([]byte(js))
		chk.Error(err)
		chk.Contains(err.Error(), "second unmarshal error")
	}
}

func TestMU_Marshal_MarshallerReturnsErrors(t *testing.T) {
	chk := assert.New(t)
	//
	type Person struct {
		jsmu.TypeName `jsmu:"person"`
		Name          string `json:"name"`
		Age           int    `json:"age"`
	}
	//
	mock := &jsmu.MockMarshaller{}
	//
	mu := jsmu.MU{
		Marshaller: mock,
	}
	mu.MustRegister(&Person{})
	//
	{
		// First call to Marshal errors
		mock.MarshalImplementation = func(v interface{}) ([]byte, error) {
			return nil, errors.Errorf("first marshal error")
		}
		_, err := mu.Marshal(&Person{})
		chk.Error(err)
		chk.Contains(err.Error(), "first marshal error")
	}
	//
	{
		// Second call to Unmarshal errors
		calls := -1
		mock.MarshalImplementation = func(v interface{}) ([]byte, error) {
			calls++
			if calls == 0 {
				return jsmu.DefaultMarshaller.Marshal(v)
			}
			return nil, errors.Errorf("second marshal error")
		}
		_, err := mu.Marshal(&Person{})
		chk.Error(err)
		chk.Contains(err.Error(), "second marshal error")
	}
	//
	// Now duplicate the above tests when v is already an jsmu.Enveloper
	envelope := jsmu.DefaultEnveloperFunc()
	envelope.SetMessage(&Person{})
	{
		// First call to Marshal errors
		mock.MarshalImplementation = func(v interface{}) ([]byte, error) {
			return nil, errors.Errorf("first marshal error")
		}
		_, err := mu.Marshal(envelope)
		chk.Error(err)
		chk.Contains(err.Error(), "first marshal error")
	}
	//
	{
		// Second call to Unmarshal errors
		calls := -1
		mock.MarshalImplementation = func(v interface{}) ([]byte, error) {
			calls++
			if calls == 0 {
				return jsmu.DefaultMarshaller.Marshal(v)
			}
			return nil, errors.Errorf("second marshal error")
		}
		_, err := mu.Marshal(envelope)
		chk.Error(err)
		chk.Contains(err.Error(), "second marshal error")
	}
}
