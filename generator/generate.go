package generator

import (
	"fmt"
	"go/importer"
	"go/token"
	"go/types"

	"github.com/FanFani4/goverter/builder"
	"github.com/FanFani4/goverter/comments"
	"github.com/FanFani4/goverter/namer"
	"github.com/FanFani4/goverter/xtype"
	"github.com/dave/jennifer/jen"
)

// Config the generate config.
type Config struct {
	Name string
}

// BuildSteps that'll used for generation.
var BuildSteps = []builder.Builder{
	&builder.BasicTargetPointerRule{},
	&builder.Pointer{},
	&builder.TargetPointer{},
	&builder.Basic{},
	&builder.Struct{},
	&builder.List{},
	&builder.Map{},
}

// Generate generates a jen.File containing converters.
func Generate(pattern string, mapping []comments.Converter, config Config) (*jen.File, error) {
	sources, err := importer.ForCompiler(token.NewFileSet(), "source", nil).Import(pattern)
	if err != nil {
		return nil, err
	}
	file := jen.NewFile(config.Name)
	file.HeaderComment("// Code generated by github.com/FanFani4/goverter, DO NOT EDIT.")

	for _, converter := range mapping {
		obj := sources.Scope().Lookup(converter.Name)
		if obj == nil {
			return nil, fmt.Errorf("%s: could not find %s", pattern, converter.Name)
		}

		// create the converter struct
		file.Type().Id(converter.Config.Name).Struct()

		gen := generator{
			namer:             namer.New(),
			file:              file,
			name:              converter.Config.Name,
			lookup:            map[xtype.Signature]*methodDefinition{},
			extend:            map[xtype.Signature]*methodDefinition{},
			ignore_unexported: converter.Config.IgnoreUnexported,
			case_insensitive:  converter.Config.CaseInsensitive,
		}

		interf := obj.Type().Underlying().(*types.Interface)

		if err := gen.parseExtend(obj.Type(), sources.Scope(), converter.Config.ExtendMethods); err != nil {
			return nil, fmt.Errorf("Error while parsing extend methods: %s", err)
		}

		// we checked in comments, that it is an interface
		for i := 0; i < interf.NumMethods(); i++ {
			method := interf.Method(i)
			converterMethod := comments.Method{}

			if m, ok := converter.Methods[method.Name()]; ok {
				converterMethod = m
			}
			if err := gen.registerMethod(method, converterMethod); err != nil {
				return nil, fmt.Errorf("Error while creating converter method:\n    %s\n\n%s", method.String(), err)
			}
		}
		if err := gen.createMethods(); err != nil {
			return nil, err
		}
	}
	return file, nil
}
