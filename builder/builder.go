package builder

import (
	"github.com/FanFani4/goverter/namer"
	"github.com/FanFani4/goverter/xtype"
	"github.com/dave/jennifer/jen"
)

// Builder builds converter implementations, and can decide if it can handle the given type.
type Builder interface {
	// Matches returns true, if the builder can create handle the given types
	Matches(source, target *xtype.Type) bool
	// Build creates conversion source code for the given source and target type.
	Build(gen Generator, ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error)
}

// Generator checks all existing builders if they can create a conversion implementations for the given source and target type
// If no one Builder#Matches then, an error is returned.
type Generator interface {
	Build(ctx *MethodContext, sourceID *xtype.JenID, source, target *xtype.Type) ([]jen.Code, *xtype.JenID, *Error)
}

// MethodContext exposes information for the current method.
type MethodContext struct {
	*namer.Namer
	Mapping         map[string]string
	IgnoredFields   map[string]struct{}
	IdentityMapping map[string]struct{}
	Signature       xtype.Signature
	PointerChange   bool
	SkipUnexported  bool
	MapLower        bool
}
