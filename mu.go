package jsmu

import (
	"encoding/json"
	"reflect"
	"sync"

	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/set"
)

// wrappedPtr is a special type indicating the value pointed-at could change storage locations as part of
// encoding.  For example if a []*T is registered and we later want to unmarshal it we need to
// pass a *[]*T into Unmarshal which is then going to change the slice allocation and affect the final
// data we return as the message.
type wrappedPtr struct {
	v   reflect.Value
	ptr interface{}
}

// MU marshals|unmarshals JavaScript encoded messages.
type MU struct {
	//
	// EnveloperFn is the function that creates a new envelope when one is needed.  If not provided
	// (aka left as nil) then DefaultEnveloperFunc is used and your JSON must conform to the format
	// described for type Envelope.
	EnveloperFn EnveloperFunc
	//
	// Marshaller is the implementation for marshalling and unmarshalling data.  If not provided
	// (aka left as nil) then a DefaultMarshaller is used; DefaultMarshaller uses json.Marshal() and
	// json.Unmarshal().
	Marshaller Marshaller
	//
	// StructTag specifies the struct tag name to use when inspecting types during register.  If not
	// set will default to "jsmu".
	StructTag string
	//
	known   sync.Map
	reverse sync.Map
}

// MustRegister is similar to Register() except it panic if an error is returned.
func (me *MU) MustRegister(value interface{}, opts ...interface{}) {
	if err := me.Register(value, opts...); err != nil {
		panic(err)
	}
}

// Register registers a type instance with the MU.
//
// Calls to Register are not goroutine safe; it is the caller's responsibility to coordinate locking if
// registering values from goroutines.
func (me *MU) Register(value interface{}, opts ...interface{}) error {
	// Use configured or default struct tag "jsmu".
	tagName := me.StructTag
	if tagName == "" {
		tagName = "jsmu"
	}
	//
	// Look for following values in our options.
	var typeName string
	var factoryFn MessageFunc
	for _, opt := range opts {
		switch o := opt.(type) {
		case MessageFunc:
			factoryFn = o
		case TypeName:
			typeName = string(o)
		}
	}
	if typeName == "" { // If not in opts then inspect the value.
		// TypeName not in opts; inspect value.
		typInfo := set.TypeCache.Stat(value)
		for _, field := range typInfo.StructFields {
			if field.Type == typeTypeName {
				typeName = field.Tag.Get(tagName)
			}
		}
	}
	if typeName == "" { // Type name was not found.
		return errors.Go(ErrMissingTypeName).Type(value)
	}
	//
	// Now we know our map key
	if _, ok := me.known.Load(typeName); ok {
		return errors.Go(ErrAlreadyRegistered).Type(value).Tag("type-name", typeName)
	}
	//
	typ := reflect.TypeOf(value)
	if factoryFn == nil {
		factoryFn = func() (interface{}, error) {
			// Save and shadow outer t.
			t := typ
			// Create a new instance and make writable; i.e. instantiate if t is a pointer.
			v := reflect.New(t)
			set.Writable(v)
			// When creating a slice we need to hang onto that initial pointer.
			if t.Kind() == reflect.Slice {
				rv := &wrappedPtr{
					v:   v,
					ptr: v.Interface(),
				}
				return rv, nil
			}
			return reflect.Indirect(v).Interface(), nil
		}
	}
	me.known.Store(typeName, factoryFn)
	me.reverse.Store(typ, typeName)
	//
	if me.EnveloperFn == nil {
		me.EnveloperFn = DefaultEnveloperFunc
	}
	if me.Marshaller == nil {
		me.Marshaller = DefaultMarshaller
	}
	//
	return nil
}

// Marshal marshals the incoming value.  If value is already an Enveloper then MU.Marshaller.Marshal(value)
// is returned.  Otherwise value is wrapped in a new envelope and MU.Marshaller.Marshal(NewEnvelope{value})
// is returned.
func (me *MU) Marshal(value interface{}) ([]byte, error) {
	if me.Marshaller == nil {
		return nil, errors.Go(ErrNothingRegistered)
	}
	//
	var reverse interface{}
	var envelope Enveloper
	var ok bool
	if envelope, ok = value.(Enveloper); ok {
		var raw json.RawMessage
		var err error
		if raw, err = me.Marshaller.Marshal(envelope.GetMessage()); err != nil {
			return nil, errors.Go(err)
		}
		envelope.SetRawMessage(raw)
		return me.Marshaller.Marshal(envelope)
	} else if reverse, ok = me.reverse.Load(reflect.TypeOf(value)); !ok {
		return nil, errors.Go(ErrUnregisteredType).Type(value)
	}
	//
	envelope = me.EnveloperFn()
	var raw json.RawMessage
	var err error
	if raw, err = me.Marshaller.Marshal(value); err != nil {
		return nil, errors.Go(err)
	}
	envelope.SetRawMessage(raw)
	envelope.SetTypeName(reverse.(string))
	//
	return me.Marshaller.Marshal(envelope)
}

// Unmarshal unmarshals a JSON string into a registered type and the envelope that contained it.
// An error is returned if the JSON can not be unmarshalled or the type is unregistered.
func (me *MU) Unmarshal(data []byte) (Enveloper, error) {
	if me.Marshaller == nil {
		return nil, errors.Go(ErrNothingRegistered)
	}
	//
	var envelope Enveloper
	var concrete, known interface{}
	var ok bool
	var err error
	//
	envelope = me.EnveloperFn()
	//
	if err = me.Marshaller.Unmarshal(data, envelope); err != nil {
		return nil, errors.Go(err)
	} else if known, ok = me.known.Load(envelope.GetTypeName()); !ok {
		return nil, errors.Go(ErrUnregisteredType).Tag("type", envelope.GetTypeName())
	}
	//
	factoryFn := known.(MessageFunc)
	if concrete, err = factoryFn(); err != nil {
		return nil, errors.Go(err)
	} else if sliceWrapper, ok := concrete.(*wrappedPtr); ok {
		if err = me.Marshaller.Unmarshal(envelope.GetRawMessage(), sliceWrapper.ptr); err != nil {
			return nil, errors.Go(err)
		}
		envelope.SetMessage(reflect.Indirect(sliceWrapper.v).Interface())
		return envelope, nil
	} else if err = me.Marshaller.Unmarshal(envelope.GetRawMessage(), concrete); err != nil {
		return nil, errors.Go(err)
	}
	envelope.SetMessage(concrete)
	//
	return envelope, nil
}
