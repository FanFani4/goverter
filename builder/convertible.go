package builder

import (
	"fmt"

	"github.com/FanFani4/goverter/xtype"
	"github.com/dave/jennifer/jen"
)

// Basic handles basic data types.
type Convertible struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Convertible) Matches(source, target *xtype.Type) bool {
	trgt := target
	if target.Pointer {
		trgt = target.PointerInner
	}

	src := source
	if src.Pointer {
		src = source.PointerInner

	}

	return src.ConvertibleTo(trgt)
}

// Build creates conversion source code for the given source and target type.
func (*Convertible) Build(_ Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {

	if !source.Pointer && !target.Pointer {
		return nil, xtype.OtherID(target.TypeAsJen().Call(sourceID.Code)), nil
	}

	stmt := []jen.Code{}

	name := ctx.Name(source.ID())
	if source.Pointer {
		ifBlock := jen.Id(name).Op("=").Id(xtype.OtherID(target.TypeAsJen().Call(jen.Op("*").Id(sourceID.Code.GoString()))).Code.GoString())
		stmt = append(stmt, jen.Var().Id(name).Id(target.TypeAsJen().GoString()))
		stmt = append(stmt, jen.If(jen.Id(sourceID.Code.GoString()).Op("!=").Nil().Block(ifBlock)))
	}

	resultID := jen.Id(name)
	if target.Pointer {
		resultID = jen.Op("&").Id(name)
	}

	for _, s := range stmt {
		fmt.Printf("%#v", s)
	}
	return stmt, xtype.OtherID(resultID), nil
}
