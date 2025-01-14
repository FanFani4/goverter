//go:generate go run github.com/FanFani4/goverter/cmd/goverter github.com/FanFani4/goverter/example/simple
package simple

// goverter:converter
type Converter interface {
	Convert(source []Input) []Output
}

type Input struct {
	Name string
	Age  int
}

type Output struct {
	Name string
	Age  int
}
