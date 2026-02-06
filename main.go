package main

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/jirwin/protoc-gen-actionqueue/internal/pgactionqueue"
)

func main() {
	minEdition := int32(descriptorpb.Edition_EDITION_PROTO2)
	maxEdition := int32(descriptorpb.Edition_EDITION_2023)
	features := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL | pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	pgs.Init(
		pgs.DebugEnv("DEBUG_PG_ACTIONQUEUE"),
		pgs.SupportedFeatures(&features),
		pgs.MinimumEdition(&minEdition),
		pgs.MaximumEdition(&maxEdition),
	).
		RegisterModule(pgactionqueue.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
