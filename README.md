<p align="center">
    <img width="300" src=".github/logo.svg" />
</p>

<h1 align="center">goverter</h1>
<p align="center"><i>a "type-safe Go converter" generator</i></p>
<p align="center">
    <a href="https://github.com/FanFani4/goverter/actions/workflows/build.yml">
        <img alt="Build Status" src="https://github.com/FanFani4/goverter/actions/workflows/build.yml/badge.svg">
    </a>
     <a href="https://codecov.io/gh/FanFani4/goverter">
        <img alt="codecov" src="https://codecov.io/gh/FanFani4/goverter/branch/main/graph/badge.svg">
    </a>
    <a href="https://goreportcard.com/report/github.com/FanFani4/goverter">
        <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/FanFani4/goverter">
    </a>
    <a href="https://pkg.go.dev/github.com/FanFani4/goverter">
        <img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/FanFani4/goverter.svg">
    </a>
    <a href="https://github.com/FanFani4/goverter/releases/latest">
        <img alt="latest release" src="https://img.shields.io/github/release/FanFani4/goverter.svg">
    </a>
</p>

goverter is a tool for creating type-safe converters. All you have to
do is create an interface and execute goverter. The project is meant as
alternative to [jinzhu/copier](https://github.com/jinzhu/copier) that doesn't
use reflection.

## Features

- Automatic conversion of builtin types
  ([`house` example](https://github.com/FanFani4/goverter/blob/main/example/house)), this includes:
  - slices, maps, named types, primitive types, pointers
  - structs with same fields
- Extend parts of the conversion with your own implementation:
  [Docs](#extend-with-custom-implementation)
- Optional return of an error: [Docs](#errors)
- Awesome error
  messages: [mismatch type test](https://github.com/FanFani4/goverter/blob/main/scenario/7_error_nested_mismatch.yml)
- No reflection in the generated code

## Usage

1. Create a go modules project if you haven't done so already
   ```bash
   $ go mod init module-name
   ```
1. Add `goverter` as dependency to your project

   ```bash
   $ go get github.com/FanFani4/goverter
   ```

1. Create your converter interface and mark it with a comment containing `goverter:converter`

   `input.go`

   ```go
   package example

   // goverter:converter
   type Converter interface {
     Convert(source []Input) []Output
   }

   type Input struct {
     Name string
     Age int
   }
   type Output struct {
     Name string
     Age int
   }
   ```

1. Run `goverter`:

   ```
   $ go run github.com/FanFani4/goverter/cmd/goverter module-name-in-full
   # example
   $ go run github.com/FanFani4/goverter/cmd/goverter github.com/FanFani4/goverter/example/simple
   ```

1. goverter created a file at `./generated/generated.go`, it may look like this:

   ```go
   package generated

   import simple "github.com/FanFani4/goverter/example/simple"

   type ConverterImpl struct{}

   func (c *ConverterImpl) Convert(source []simple.Input) []simple.Output {
     simpleOutputList := make([]simple.Output, len(source))
     for i := 0; i < len(source); i++ {
       simpleOutputList[i] = c.simpleInputToSimpleOutput(source[i])
     }
     return simpleOutputList
   }
   func (c *ConverterImpl) simpleInputToSimpleOutput(source simple.Input) simple.Output {
     var simpleOutput simple.Output
     simpleOutput.Name = source.Name
     simpleOutput.Age = source.Age
     return simpleOutput
   }
   ```

## Docs

### Rename converter

With `goverter:name` you can set the name of the generated converter struct.

`input.go`

```go
// goverter:converter
// goverter:name RenamedConverter
type BadlyNamed interface {
    // .. methods
}
```

`output.go`

```go
type RenamedConverter struct {}

func (c *RenamedConverter) ...
```

### Extend with custom implementation

With `goverter:extend` you can instruct goverter to use an implementation of
your own. You can pass multiple function names to `goverter:extend`, or define
the tag multiple times.

See [`house` example sql.NullString](https://github.com/FanFani4/goverter/blob/main/example/house/input.go#L9)

`input.go`

```go
// goverter:converter
// goverter:extend IntToString
type Converter interface {
    Convert(Input) Output
}
type Input struct {Value int}
type Output struct {Value string}

// You must atleast define a source and target type. Meaning one parameter and one return.
// You can use any type you want, like struct, maps and so on.
func IntToString(i int) string {
    return fmt.Sprint(i)
}
```

#### Reuse generated converter

If you need access to the generated converter, you can define it as first
parameter.

```go
func IntToString(c Converter, i int) string {
    // c.DoSomething()..
    return fmt.Sprint(i)
}
```

#### Errors

Sometimes, custom conversion may fail, in this case goverter allows you to
define a second return parameter which must be type `error`.

```go
// goverter:converter
// goverter:extend IntToString
type Converter interface {
    Convert(Input) (Output, error)
}

type Input struct {Value int}
type Output struct {Value string}

func IntToString(i int) (string, error) {
    if i == 0 {
        return "", errors.New("zero is not allowed")
    }
    return fmt.Sprint(i)
}
```

_Note_: If you do this, methods on the interface that'll use this custom
implementation, must also return error as second return.

### Struct field mapping

With `goverter:map` you can map fields on structs that have the same type.

`goverter:map` takes 2 parameters.

1. source field path (fields are separated by `.`)
1. target field name

```go
// goverter:converter
type Converter interface {
    // goverter:map Name FullName
    // goverter:map Nested.Age Age
    Convert(source Input) Output
}

type Input struct {
    Name string
    Nested NestedInput
}
type NestedInput struct {
    Age int
}
type Output struct {
    FullName string
    Age int
}
```


#### Struct identity mapping

With `goverter:mapIdentity` you can instruct goverter to use the source struct
as source for the conversion to the target property.

`goverter:mapIdentity` takes multiple field names separated by space ` `.

```go
// goverter:converter
type Converter interface {
    // goverter:mapIdentity Address
    ConvertPerson(source Person) APIPerson
}

type Person struct {
    Name   string
    Street string
    City   string
}

type APIPerson struct {
    Name    string
    Address APIAddress
}

type APIAddress struct {
    Street string
    City   string
}
```

In the example goverter will fill `Address` by creating a converter from
`Person` to `APIAddress`. Example generated code:

```go
func (c *ConverterImpl) ConvertPerson(source execution.Person) execution.APIPerson {
    var structsAPIPerson execution.APIPerson
    structsAPIPerson.Name = source.Name
    structsAPIPerson.Address = c.structsPersonToStructsAPIAddress(source)
    return structsAPIPerson
}
func (c *ConverterImpl) structsPersonToStructsAPIAddress(source execution.Person) execution.APIAddress {
    var structsAPIAddress execution.APIAddress
    structsAPIAddress.Street = source.Street
    structsAPIAddress.City = source.City
    return structsAPIAddress
}
```

### Struct ignore field

With `goverter:ignore` you can ignore fields on the target struct

`goverter:ignore` takes multiple field names separated by space ` `.

```go
// goverter:converter
type Converter interface {
    // goverter:ignore Age
    Convert(source Input) Output
}

type Input struct {
    Name string
}
type Output struct {
    Name string
    Age int
}
```

## Versioning

goverter use [SemVer](http://semver.org/) for versioning the cli.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE)
file for details

_Logo by [MariaLetta](https://github.com/MariaLetta)_
