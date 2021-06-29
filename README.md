[![Documentation](https://godoc.org/github.com/nofeaturesonlybugs/jsmu?status.svg)](http://godoc.org/github.com/nofeaturesonlybugs/jsmu)
[![Go Report Card](https://goreportcard.com/badge/github.com/nofeaturesonlybugs/jsmu)](https://goreportcard.com/report/github.com/nofeaturesonlybugs/jsmu)
[![Build Status](https://travis-ci.com/nofeaturesonlybugs/jsmu.svg?branch=master)](https://travis-ci.com/nofeaturesonlybugs/jsmu)
[![codecov](https://codecov.io/gh/nofeaturesonlybugs/jsmu/branch/master/graph/badge.svg)](https://codecov.io/gh/nofeaturesonlybugs/jsmu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`jsmu` aka `JavaScript Marshal|Unmarshal`.  Pronounce as `j◦s◦mew`.

## `jsmu.MU`  
`jsmu.MU` abstracts the typical 2-pass JSON marshal|unmarshal into a single pass.  

Create an instance of `jsmu.MU`:  
```go
mu := &jsmu.MU{}
```

Register types that will be JSON marshalled or unmarshalled:  
```go
mu.Register( &SomeType{} )
mu.Register( &OtherType{} )
```

To unmarshal:  
```go
var envelope jsmu.Enveloper
var err error
for data := range jsonCh { // jsonCh is some stream returning blobs of JSON
    if envelope, err = mu.Unmarshal([]byte(data)); err != nil {
        // YOUSHOULD Handle err.
        continue // break, goto, return, etc.
    }
    switch message := envelope.GetMessage().(type) {
        case *SomeType:
        case *OtherType:
    }
}
```

To marshal:  
```go
var buf []byte
var err error
some, other := &SomeType{}, &OtherType{}
if buf, err = mu.Marshal(some); err != nil {
    // YOUSHOULD Handle err.
}
// YOUSHOULD Do something with buf.
if buf, err = mu.Marshal(other); err != nil {
    // YOUSHOULD Handle err
}
// YOUSHOULD Do something with buf.
```

### Envelopes  
`jsmu` expects JSON encoded data to be wrapped inside an envelope.  This envelope must contain some type name information (data type `string`) and the message payload.  The default envelope structure is:  
```js
{
    // string, required : Must uniquely identify the contents of "message".
    "type"    : "",
    // string, optional : Provided as a convenience for round trip message flow.
    "id"      : "",
    // The data.
    "message" : any
}
```

Envelopes are created by the `EnveloperFn` member of `jsmu.MU`; if you don't provide a value for `EnveloperFn` then `jsmu.DefaultEnveloperFunc` is used and JSON is expected to conform to the above structure.  

If you want to use a different JSON structure you can implement the `jsmu.Enveloper` interface and provide an appropriate `EnveloperFn` value when instantiating your `jsmu.MU`; see example `MU (CustomEnvelope)`.

### The Dual Nature of `MU.Marshal()`  
`MU.Marshal()` *always* returns JSON representing a message within an envelope.  
* Calling `MU.Marshal(value)` when `value` is already a `jsmu.Enveloper` yields:
    1. `value` is marshalled as-is.

* Calling `MU.Marshal(value)` when `value` *is not* a `jsmu.Enveloper` yields:
    1. A new envelope is created with `MU.EnveloperFn`.
    2. `value` is placed in the envelope.
    3. The newly created envelope is marshalled.
