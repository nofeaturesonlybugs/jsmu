[![Go Reference](https://pkg.go.dev/badge/github.com/nofeaturesonlybugs/jsmu.svg)](https://pkg.go.dev/github.com/nofeaturesonlybugs/jsmu)
[![Go Report Card](https://goreportcard.com/badge/github.com/nofeaturesonlybugs/jsmu)](https://goreportcard.com/report/github.com/nofeaturesonlybugs/jsmu)
[![Build Status](https://app.travis-ci.com/nofeaturesonlybugs/jsmu.svg?branch=master)](https://app.travis-ci.com/nofeaturesonlybugs/jsmu)
[![codecov](https://codecov.io/gh/nofeaturesonlybugs/jsmu/branch/master/graph/badge.svg)](https://codecov.io/gh/nofeaturesonlybugs/jsmu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`jsmu` aka `JavaScript Marshal|Unmarshal`. Pronounce as `j◦s◦mew`.

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

`jsmu` expects JSON encoded data to be wrapped inside an envelope. This envelope must contain some type name information (data type `string`) and the message payload. The default envelope structure is:

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

`MU.Marshal()` _always_ returns JSON representing a message within an envelope.

-   Calling `MU.Marshal(value)` when `value` is already a `jsmu.Enveloper` yields:

    1. `value` is marshalled as-is.

-   Calling `MU.Marshal(value)` when `value` _is not_ a `jsmu.Enveloper` yields:
    1. A new envelope is created with `MU.EnveloperFn`.
    2. `value` is placed in the envelope.
    3. The newly created envelope is marshalled.

## Performance

`jsmu.MU` has not yet been optimized. However the extra overhead is not much slower than handling 2-pass JSON directly in your application (at least for the simple test I put together):

```
goos: windows
goarch: amd64
pkg: github.com/nofeaturesonlybugs/jsmu
cpu: Intel(R) Core(TM) i7-7700K CPU @ 4.20GHz
BenchmarkMU/2-pass_unmarshal_limit_5-8             52425             23602 ns/op            7488 B/op        140 allocs/op
BenchmarkMU/jsmu_unmarshal_limit_5-8               50223             24460 ns/op            7008 B/op        130 allocs/op
BenchmarkMU/2-pass_unmarshal_limit_100-8            2731            466576 ns/op          147256 B/op       2800 allocs/op
BenchmarkMU/jsmu_unmarshal_limit_100-8              2668            476442 ns/op          137656 B/op       2600 allocs/op
BenchmarkMU/2-pass_unmarshal_limit_250-8             994           1171133 ns/op          367593 B/op       7000 allocs/op
BenchmarkMU/jsmu_unmarshal_limit_250-8              1014           1194731 ns/op          343593 B/op       6500 allocs/op
BenchmarkMU/2-pass_unmarshal_limit_500-8             519           2343670 ns/op          734043 B/op      14000 allocs/op
BenchmarkMU/jsmu_unmarshal_limit_500-8               513           2394790 ns/op          686043 B/op      13000 allocs/op
BenchmarkMU/2-pass_unmarshal_limit_1000-8            255           4734993 ns/op         1467949 B/op      28000 allocs/op
BenchmarkMU/jsmu_unmarshal_limit_1000-8              244           4788199 ns/op         1371949 B/op      26000 allocs/op
PASS
ok      github.com/nofeaturesonlybugs/jsmu      14.532s
```
