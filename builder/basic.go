package builder

import (
	"github.com/FanFani4/goverter/xtype"
	"github.com/dave/jennifer/jen"
)

// Basic handles basic data types.
type Basic struct{}

// Matches returns true, if the builder can create handle the given types.
func (*Basic) Matches(source, target *xtype.Type) bool {
	return source.Basic && target.Basic &&
		source.BasicType.Kind() == target.BasicType.Kind()
}

// Build creates conversion source code for the given source and target type.
func (*Basic) Build(_ Generator, _ *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	if target.Named || (!target.Named && source.Named) {
		return nil, xtype.OtherID(target.TypeAsJen().Call(sourceID.Code)), nil
	}
	return nil, sourceID, nil
}

// BasicTargetPointerRule handles edge conditions if the target type is a pointer.
type BasicTargetPointerRule struct{}

// Matches returns true, if the builder can create handle the given types.
func (*BasicTargetPointerRule) Matches(source, target *xtype.Type) bool {
	return source.Basic && target.Pointer && target.PointerInner.Basic
}

// Build creates conversion source code for the given source and target type.
func (*BasicTargetPointerRule) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	name := ctx.Name(target.ID())

	stmt, id, err := gen.Build(ctx, sourceID, source, target.PointerInner)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.T.String(),
			TargetID:   "*",
			TargetType: target.PointerInner.T.String(),
		})
	}
	stmt = append(stmt, jen.Id(name).Op(":=").Add(id.Code))
	newID := jen.Op("&").Id(name)

	return stmt, xtype.OtherID(newID), err
}

type BasicSourcePtrRule struct{}

// Matches returns true, if the builder can create handle the given types.
func (*BasicSourcePtrRule) Matches(source, target *xtype.Type) bool {
	return source.Pointer && source.PointerInner.Basic && target.Basic && source.PointerInner.BasicType.Kind() == target.BasicType.Kind()
}

func (*BasicSourcePtrRule) Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error) {
	name := ctx.Name(target.ID())

	stmt, _, err := gen.Build(ctx, sourceID, source.PointerInner, target)
	if err != nil {
		return nil, nil, err.Lift(&Path{
			SourceID:   "*",
			SourceType: source.PointerInner.T.String(),
			TargetID:   "*",
			TargetType: target.T.String(),
		})
	}

	src := sourceID.Code.GoString()
	stmt = append(stmt, jen.Var().Id(name).Id(target.T.String()))
	stmt = append(stmt, jen.If(sourceID.Code.Op("!=").Nil().Block(jen.Id(name).Op("=").Op("*").Id(src))))
	newID := jen.Id(name)

	return stmt, xtype.OtherID(newID), err
}
