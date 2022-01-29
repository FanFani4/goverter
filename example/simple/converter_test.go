package simple_test

import (
	"testing"

	"github.com/FanFani4/goverter/example/simple"
	"github.com/FanFani4/goverter/example/simple/generated"
	"github.com/stretchr/testify/require"
)

func TestConverter(t *testing.T) {
	var c simple.Converter = &generated.ConverterImpl{}
	inputs := []simple.Input{
		{Age: 5, Name: "FanFani4"},
		{Age: 75, Name: "other"},
	}

	actual := c.Convert(inputs)

	expected := []simple.Output{
		{Age: 5, Name: "FanFani4"},
		{Age: 75, Name: "other"},
	}

	require.Equal(t, expected, actual)
}
